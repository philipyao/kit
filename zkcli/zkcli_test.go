package zkcli

import (
    //"time"
    //"fmt"
    "testing"
    //"github.com/samuel/go-zookeeper/zk"
)

//func TestWatch(t *testing.T) {
//    zkConn, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
//    if err != nil {
//        t.Fatalf("Connect returned error: %+v", err)
//    }
//    defer zkConn.Close()
//
//    if err := zkConn.Conn().Delete("/gozk-test-w1", -1); err != nil && err != zk.ErrNoNode {
//        t.Fatalf("Delete returned error: %+v", err)
//    }
//
//    testPath, err := zkConn.Conn().Create("/gozk-test-w1", []byte{}, 0, zk.WorldACL(zk.PermAll))
//    if err != nil {
//        t.Fatalf("Create returned: %+v", err)
//    }
//    stop := make(chan struct{}, 1)
//    err = zkConn.Watch(testPath, func(p string, d []byte, e error){
//        fmt.Println("w1: ", p, d, e)
//    }, stop)
//    if err != nil {
//        t.Fatal(err)
//    }
//    stop2 := make(chan struct{}, 1)
//    err = zkConn.Watch("/notexist", func(p string, d []byte, e error){
//        fmt.Println("w2: ", p, d, e)
//    }, stop2)
//    if err != nil {
//        t.Fatal(err)
//    }
//
//    time.Sleep(time.Second * 60)
//    stop <- struct{}{}
//    stop2 <- struct{}{}
//    time.Sleep(time.Second)
//}

//func TestWatchNode(t *testing.T) {
//    zkConn, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
//    if err != nil {
//        t.Fatalf("Connect returned error: %+v", err)
//    }
//    defer zkConn.Close()
//
//    zkPath := "/gozk-test-w11"
//    err = zkConn.WriteEphemeral(zkPath, []byte{})
//    if err != nil {
//        t.Fatalf("Create returned: %+v", err)
//    }
//    exit := make(chan struct{})
//    events := make(chan *EventDataNode, 10)
//    zkConn.WatchNode(zkPath, events, exit)
//
//    go func() {
//        time.Sleep(time.Second * 1)
//        zkConn2, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
//        if err != nil {
//            t.Fatalf("Connect returned error: %v", err)
//        }
//        defer zkConn2.Close()
//
//        time.Sleep(time.Second * 1)
//        err = zkConn2.WriteEphemeral(zkPath, []byte("hello"))
//        if err != nil {
//            t.Fatalf("WriteEphemeral returned: %v", err)
//        }
//        time.Sleep(time.Second * 1)
//        err = zkConn2.WriteEphemeral(zkPath, []byte("world"))
//        if err != nil {
//            t.Fatalf("WriteEphemeral returned: %v", err)
//        }
//    }()
//
//    go func() {
//        for ev := range events {
//            fmt.Printf("received event: %+v, time %v\n", ev, time.Now().Unix())
//        }
//        t.Logf("receive goroutine exit.")
//    }()
//
//    <-time.After(8 * time.Second)
//    exit <- struct{}{}
//    close(events)
//
//    time.Sleep(1 * time.Second)
//}

func TestMakeDirP(t *testing.T) {
    zkConn, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
    if err != nil {
        t.Fatalf("Connect returned error: %+v", err)
    }
    defer zkConn.Close()

    zkPath := "/gozk-test-dir/foo/bar/"
    err = zkConn.MakeDirP(zkPath)
    if err != nil {
        t.Fatalf("MakeDirP returned: %v", err)
    }
}

//func TestWatchDir(t *testing.T) {
//    zkConn, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
//    if err != nil {
//        t.Fatalf("Connect returned error: %+v", err)
//    }
//    defer zkConn.Close()
//
//    zkPath := "/gozk-test-w12"
//    err = zkConn.Write(zkPath, []byte{})
//    if err != nil {
//        t.Fatalf("Create returned: %+v", err)
//    }
//    exit := make(chan struct{})
//    events := make(chan *EventDataChild, 10)
//    zkConn.WatchDir(zkPath, events, exit)
//
//    go func() {
//        time.Sleep(time.Second * 2)
//        zkConn2, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
//        if err != nil {
//            t.Fatalf("Connect returned error: %v", err)
//        }
//        defer zkConn2.Close()
//
//        time.Sleep(time.Second * 2)
//        err = zkConn2.WriteEphemeral(zkPath + "/service1", []byte("1"))
//        if err != nil {
//            t.Fatalf("WriteEphemeral returned: %v", err)
//        }
//        err = zkConn2.WriteEphemeral(zkPath + "/service2", []byte("2"))
//        if err != nil {
//            t.Fatalf("WriteEphemeral returned: %v", err)
//        }
//
//        time.Sleep(time.Second * 2)
//    }()
//
//    go func() {
//        for ev := range events {
//            fmt.Printf("received event: %+v, time %v\n", ev, time.Now().Unix())
//        }
//        t.Logf("receive goroutine exit.")
//    }()
//
//    <-time.After(8 * time.Second)
//    exit <- struct{}{}
//    close(events)
//
//    time.Sleep(1 * time.Second)
//}

//func TestWatchChildren(t *testing.T) {
//    zkConn, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
//    if err != nil {
//        t.Fatalf("Connect returned error: %+v", err)
//    }
//    defer zkConn.Conn().Close()
//
//    if err := zkConn.Conn().Delete("/gozk-test-wc", -1); err != nil && err != zk.ErrNoNode {
//        t.Fatalf("Delete returned error: %+v", err)
//    }
//    testPath, err := zkConn.Conn().Create("/gozk-test-wc", []byte{}, 0, zk.WorldACL(zk.PermAll))
//    if err != nil {
//        t.Fatalf("Create returned: %+v", err)
//    }
//    stop := make(chan struct{}, 1)
//    err = zkConn.WatchChildren(testPath, func(p string, c []string, e error){
//        fmt.Println("wc: ", p, c, e)
//    }, stop)
//    if err != nil {
//        t.Fatal(err)
//    }
//    time.Sleep(time.Second * 120)
//    stop <- struct{}{}
//    time.Sleep(time.Second)
//}
//
//func TestCreateEphemeral(t *testing.T) {
//    zkConn, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
//    if err != nil {
//        t.Fatalf("Connect returned error: %+v", err)
//    }
//    defer zkConn.Conn().Close()
//
//    err = zkConn.CreateEphemeral("/go-zktest-ephemeral1", []byte{})
//    if err != nil {
//        t.Fatal(err)
//    }
//    //模拟session断开
//    //zkConn.Conn().TmpCloseConn()
//
//    //time.Sleep(60 * time.Second)
//}
//
//func TestGetChildren(t *testing.T) {
//    zkConn, err := Connect("10.1.164.20:2181,10.1.164.20:2182")
//    if err != nil {
//        t.Fatalf("Connect returned error: %+v", err)
//    }
//    defer zkConn.Conn().Close()
//
//    err = zkConn.Create("/go-zktest-children", []byte{})
//    if err != nil {
//        t.Fatal(err)
//    }
//
//    for i := 0; i < 5; i++ {
//        err = zkConn.Create(fmt.Sprintf("/go-zktest-children/c%v", i+1), []byte(fmt.Sprintf("val%v", i+1)))
//        if err != nil {
//            t.Fatal(err)
//        }
//    }
//
//    results, err := zkConn.GetChildren("/go-zktest-children")
//    if err != nil {
//        t.Fatal(err)
//    }
//    for k, v := range results {
//        t.Logf("key %v, value %v", k, string(v))
//    }
//
//}
