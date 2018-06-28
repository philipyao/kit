package config

import (
    "gopkg.in/yaml.v2"
)

var cfgy ConfigYaml

type ConfigYaml struct {
}

func (cy ConfigYaml) Parse(data []byte, out interface{}) error {
    return yaml.Unmarshal(data, out)
}

func init() {
    Register(CONFIG_TP_YAML, cfgy)
}