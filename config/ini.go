package config

import (
    ini "github.com/sspencer/go-ini"
)

var cfgi ConfigIni

type ConfigIni struct {
}

func (ci ConfigIni) Parse(data []byte, out interface{}) error {
    return ini.Unmarshal(data, out)
}


func init() {
    Register(CONFIG_TP_INI, cfgi)
}