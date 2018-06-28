package zkcli

import (
    "errors"
    "fmt"
    "log"
    "path/filepath"
    "strings"
    "time"

    "github.com/samuel/go-zookeeper/zk"
)

//回调：监视节点变化
type FuncWatchCallback func(path string, data []byte, err error)

//回调：监视子节点改变
type FuncWatchChildrenCallback func(path string, children []string, err error)

const (
    //zk session的超时时间，在此时间内，session可以自动保活
    //如果在此时间内，与特定server的连接断开，自动尝试重连其他服务器
    //重连成功之后，临时节点和watch依然有效
    DefaultConnectTimeout = 5
)

type EventDataChild struct {
    Closed bool
    Err    error
    Adds   map[string]string //新增path/value
    Dels   []string          //删除path
}
type EventDataNode struct {
    Closed bool
    Err    error
    Path   string
    Value  string
    OldVal string
}

type Conn struct {
    conn *zk.Conn
}

func (c *Conn) Conn() *zk.Conn {
    return c.conn
}
func (c *Conn) SetConn(conn *zk.Conn) {
    c.conn = conn
}

func (c *Conn) Close() {
    c.conn.Close()
}

func (c *Conn) Exists(path string) (bool, error) {
    exist, _, err := c.Conn().Exists(path)
    if err != nil {
        return false, err
    }
    return exist, nil
}

func (c *Conn) Get(path string) ([]byte, error) {
    data, _, err := c.Conn().Get(path)
    return data, err
}

func (c *Conn) GetChildren(path string) (map[string][]byte, error) {
    children := make(map[string][]byte)
    nodes, _, err := c.Conn().Children(path)
    if err != nil {
        return nil, err
    }
    for _, child := range nodes {
        tp := path + "/" + child
        data, err := c.Get(tp)
        if err != nil {
            return nil, fmt.Errorf("child node %v get error %v", tp, err)
        }
        children[tp] = data
    }
    return children, nil
}

func (c *Conn) Write(path string, data []byte) error {
    //永久节点
    return c.write(path, data, int32(0))
}

func (c *Conn) WriteEphemeral(path string, data []byte) error {
    //临时节点
    return c.write(path, data, int32(zk.FlagEphemeral))
}

func (c *Conn) write(path string, data []byte, flags int32) error {
    exist, stat, err := c.conn.Exists(path)
    if err != nil {
        return err
    }
    if exist {
        _, err = c.conn.Set(path, data, stat.Version)
    } else {
        //不存在则创建
        _, err = doCreate(c.conn, path, data, flags)
    }
    return err
}

func (c *Conn) Watch(path string, cb FuncWatchCallback, stopCh chan struct{}) error {
    _, ch, err := getW(c.conn, path)
    if err != nil {
        return err
    }
    go func() {
        var data []byte
        for {
            select {
            case <-stopCh:
                return
            case ev := <-ch:
                if ev.Err != nil {
                    //错误回调
                    cb(path, nil, ev.Err)
                    return
                }
                if ev.Path != path {
                    cb(path, nil, fmt.Errorf("mismatched path %v %v", ev.Path, path))
                    return
                }
            }
            // 获取变化后的节点数据
            // 并更新watcher（zookeeper的watcher是一次性的）
            data, ch, err = getW(c.conn, path)
            if err != nil {
                //错误回调
                cb(path, nil, err)
                return
            }
            //数据回调
            cb(path, data, nil)
        }
    }()

    return nil
}

func (c *Conn) Create(path string, data []byte) error {
    exist, _, err := c.Conn().Exists(path)
    if err != nil {
        return err
    }
    if exist {
        return nil
    }
    _, err = doCreate(c.Conn(), path, data, 0)
    return err
}

func (c *Conn) Delete(path string) error {
    return c.Conn().Delete(path, 0)
}

//依次创建各级目录，类似unix系统 mkdir -p
func (c *Conn) MakeDirP(dir string) error {
    if len(dir) == 0 {
        return errors.New("empty dir")
    }
    if !strings.HasPrefix(dir, "/") {
        return errors.New("invalid dir: no root")
    }
    if dir == "/" {
        return nil
    }
    dir = strings.TrimSuffix(dir, "/")

    var (
        path string
        err  error
    )
    for _, v := range strings.Split(dir, "/")[1:] {
        path = path + "/" + v
        err = c.Create(path, []byte{})
        if err != nil {
            return err
        }
    }
    return nil
}

func (c *Conn) CreateEphemeral(path string, data []byte) error {
    exist, _, err := c.Conn().Exists(path)
    if err != nil {
        return err
    }
    if exist {
        return nil
    }
    // 临时节点
    _, err = doCreate(c.Conn(), path, data, int32(zk.FlagEphemeral))
    return err
}

func (c *Conn) CreateSequence(path string, data []byte) (string, error) {
    flags := int32(zk.FlagSequence | zk.FlagEphemeral)
    return doCreate(c.Conn(), path, data, flags)
}

func (c *Conn) WatchChildren(path string, cb FuncWatchChildrenCallback, stopCh chan struct{}) error {
    _, ch, err := childrenW(c.Conn(), path)
    if err != nil {
        return err
    }
    go func() {
        for {
            select {
            case <-stopCh:
                return
            case ev := <-ch:
                if ev.Err != nil {
                    //错误回调
                    cb(path, nil, ev.Err)
                    return
                }
                if ev.Path != path {
                    cb(path, nil, fmt.Errorf("mismatched path %v %v", ev.Path, path))
                    return
                }
            }
            // 获取变化后的节点数据
            // 并更新watcher（zookeeper的watcher是一次性的）
            var children []string
            children, ch, err = childrenW(c.Conn(), path)
            if err != nil {
                //错误回调
                cb(path, nil, err)
                return
            }
            //数据回调
            cb(path, children, nil)
        }
    }()

    return nil
}

//watch一个dir的子节点，关注节点的增删
func (c *Conn) WatchDir(path string, events chan *EventDataChild, exit chan struct{}) {
    if len(path) == 0 {
        sendChildEventError(events, errors.New("empty path"))
        return
    }
    if path != "/" {
        path = strings.TrimSuffix(path, "/")
    }

    children, ch, err := childrenW(c.Conn(), path)
    if err != nil {
        sendChildEventError(events, fmt.Errorf("childrenW<%v> err: %v", path, err))
        return
    }
    //fmt.Printf("zk: children: %+v\n", children)
    var handle bool
    go func() {
        for {
            handle = true
            select {
            case <-exit:
                //上层关闭监听
                fmt.Println("zk watchDir exit by caller")
                close(events)
                return
            case ev := <-ch:
                log.Printf("watch dir trigger: %+v", ev)
                if ev.Type == zk.EventNotWatching {
                    return
                }
                if ev.Err != nil {
                    sendChildEventError(events, fmt.Errorf("watch path %v event err %v", path, ev.Err))
                    return
                }
                if ev.Path != path {
                    sendChildEventError(events, errors.New("watch path mismatch"))
                    return
                }
                if ev.Type != zk.EventNodeChildrenChanged {
                    handle = false
                }
            }
            // 获取变化后的节点数据
            // 并更新watcher（zookeeper的watcher是一次性的）
            oldChildren := children
            children, ch, err = childrenW(c.Conn(), path)
            if err != nil {
                sendChildEventError(events, fmt.Errorf("childrenW<%v> err: %v", path, err))
                return
            }
            fmt.Printf("zk: handle %v, children: %+v, old: %+v\n", handle, children, oldChildren)
            if handle {
                sendDiffChildren(c.Conn(), path, events, oldChildren, children)
            }
        }
    }()
}

func (c *Conn) WatchNode(path string, events chan *EventDataNode, exit chan struct{}) {
    if len(path) == 0 {
        sendNodeEventError(events, errors.New("empty path"))
        return
    }
    if path != "/" {
        path = strings.TrimSuffix(path, "/")
    }
    data, ch, err := getW(c.conn, path)
    if err != nil {
        sendNodeEventError(events, fmt.Errorf("getW<%v> err %v", path, err))
        return
    }
    go func() {
        for {
            select {
            case <-exit:
                //上层关闭监听
                fmt.Println("zk watchNode exit by caller")
                close(events)
                return
            case ev := <-ch:
                log.Printf("watch node trigger: %+v", ev)
                if ev.Type == zk.EventNotWatching {
                    sendNodeEventClosed(events)
                    return
                }
                if ev.Err != nil {
                    sendNodeEventError(events, fmt.Errorf("watch path %v event err %v", path, ev.Err))
                    return
                }
                if ev.Path != path {
                    sendNodeEventError(events, errors.New("watch path mismatch"))
                    return
                }
                if ev.Type == zk.EventNodeDeleted {
                    return
                }
            }
            // 获取变化后的节点数据
            // 并更新watcher（zookeeper的watcher是一次性的）
            oldData := data
            data, ch, err = getW(c.conn, path)
            if err != nil {
                sendNodeEventError(events, fmt.Errorf("getW<%v> err %v", path, err))
                return
            }
            sendDiffNode(path, events, string(oldData), string(data))
        }
    }()
}

func Connect(zkAddr string) (*Conn, error) {
    if len(zkAddr) == 0 {
        return nil, errors.New("empty zkAddr")
    }
    zks := strings.Split(zkAddr, ",")
    conn, _, err := zk.Connect(zks, time.Second*DefaultConnectTimeout)
    if err != nil {
        return nil, fmt.Errorf("err connect to zk<%v>: %v", zkAddr, err)
    }

    c := new(Conn)
    c.SetConn(conn)
    return c, nil
}

////////////////////////////////////////////////////////////////////////////

func getW(zkConn *zk.Conn, path string) ([]byte, <-chan zk.Event, error) {
    data, _, ch, err := zkConn.GetW(path)
    return data, ch, err
}
func childrenW(zkConn *zk.Conn, path string) ([]string, <-chan zk.Event, error) {
    children, _, ch, err := zkConn.ChildrenW(path)
    return children, ch, err
}

func children(zkConn *zk.Conn, path string) ([]string, error) {
    children, _, err := zkConn.Children(path)
    return children, err
}

func doCreate(zkConn *zk.Conn, path string, data []byte, flags int32) (string, error) {
    //TODO 权限控制
    acl := zk.WorldACL(zk.PermAll)
    path, err := zkConn.Create(path, data, flags, acl)
    if err != nil {
        return "", err
    }

    return path, nil
}

func sendDiffChildren(zkConn *zk.Conn, path string, events chan *EventDataChild, oldChildren, children []string) {
    var found bool
    diff := &EventDataChild{
        Adds: make(map[string]string),
    }
    fmt.Printf("sendDiffChildren: old %+v, children %+v\n", oldChildren, children)
    for _, child := range children {
        fmt.Printf("scan children to find new: curr<%v>\n", child)
        found = false
        for _, oc := range oldChildren {
            if child == oc {
                found = true
                break
            }
        }
        fmt.Printf("found: %v\n", found)
        if !found {
            cPath := filepath.Join(path, child)
            value, _, err := zkConn.Get(cPath)
            if err != nil {
                fmt.Printf("zkConn.Get<%v> err: %v\n", cPath, err)
                continue
            }
            fmt.Printf("add new<%v> to new\n", cPath)
            diff.Adds[cPath] = string(value)
        }
    }
    for _, oc := range oldChildren {
        found = false
        for _, child := range children {
            if child == oc {
                found = true
                break
            }
        }
        if !found {
            cPath := filepath.Join(path, oc)
            diff.Dels = append(diff.Dels, cPath)
        }
    }
    select {
    case events <- diff:
    default:
        //todo log
    }
}

func sendDiffNode(path string, events chan *EventDataNode, oldData, data string) {
    diff := &EventDataNode{
        Path:   path,
        Value:  data,
        OldVal: oldData,
    }
    select {
    case events <- diff:
    default:
        //todo log
    }
}

func sendChildEventError(events chan *EventDataChild, err error) {
    errData := &EventDataChild{
        Err: err,
    }
    select {
    case events <- errData:
    default:
        //todo log
    }
}

func sendChildEventClosed(events chan *EventDataChild) {
    select {
    case events <- &EventDataChild{Closed: true}:
    default:
        //todo log
    }
}

func sendNodeEventError(events chan *EventDataNode, err error) {
    errData := &EventDataNode{
        Err: err,
    }
    select {
    case events <- errData:
    default:
        //todo log
    }
}

func sendNodeEventClosed(events chan *EventDataNode) {
    select {
    case events <- &EventDataNode{Closed: true}:
    default:
        //todo log
    }
}
