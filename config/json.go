package config

import (
    "encoding/json"
)

var cfgj ConfigJson

type ConfigJson struct {
}

func (cj ConfigJson) Parse(data []byte, out interface{}) error {
    return json.Unmarshal(data, out)
}

func init() {
    Register(CONFIG_TP_JSON, cfgj)
}