package config

import (
    "fmt"
    "io/ioutil"
)

const (
    CONFIG_TP_INI       = "ini"
    CONFIG_TP_JSON      = "json"
    CONFIG_TP_YAML      = "yaml"
)

var (
    mgrs map[string]CfgManager = make(map[string]CfgManager)
)

type CfgManager interface {
    Parse(data []byte, out interface{}) error
}

func LoadConfig(cfgtp, filename string, out interface{}) error {
    if mgrs == nil {
        panic("nil mgrs!")
    }
    mgr, exist := mgrs[cfgtp]
    if !exist {
        return fmt.Errorf("cfgtp %v not support", cfgtp)
    }

    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }
    return mgr.Parse(data, out)
}

func Register(cfgtp string, mgr CfgManager) {
    mgrs[cfgtp] = mgr
}