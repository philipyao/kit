package config

import (
    "fmt"
    "testing"
)

type TCIni struct {
    Beforeyou int `ini:"before_you"`
    PdoMysql struct {
        CacheSize     int `ini:"cache_size"`
        DefaultSocket string `ini:"default_socket"`
    } `ini:"[Pdo_myqsl]"`
    Mysql struct {
        DefaultSocket string `ini:"default_socket"`
    } `ini:"[Myqsl]"`
}

type TCJson struct {
    First string   `json:"first"`
    Second struct {
        Sec1 string `json:"sec1"`
        Sec2 []uint32 `json:"sec2"`
    } `json:"second"`
}

type TCYaml struct {
    A string `yaml:"a"`
    B struct {
        RenamedC int   `yaml:"c"`
        D        []int `yaml:"d,flow"`
    } `yaml:"b,flow"`
}

func TestIni(t *testing.T) {
    var (
        err error
        filename string
    )

    // test ini
    filename = "./name1.ini"
    ti := &TCIni{}
    err = LoadConfig(CONFIG_TP_INI, filename, ti)
    if err != nil {
        t.Fatal("load ini file error %v\n", err)
    }
    fmt.Printf("ini struct: %+v\n", ti)
}


func TestJson(t *testing.T) {
    var (
        err error
        filename string
    )

    // test json
    filename = "./name2.json"
    tj := &TCJson{}
    err = LoadConfig(CONFIG_TP_JSON, filename, tj)
    if err != nil {
        t.Error("load json file error %v\n", err)
    }
    fmt.Printf("json struct: %+v\n", tj)
}


func TestYaml(t *testing.T) {
    var (
        err error
        filename string
    )

    // test yaml
    filename = "./name3.yaml"
    ty := &TCYaml{}
    err = LoadConfig(CONFIG_TP_YAML, filename, ty)
    if err != nil {
        t.Fatal("load yaml file error %v\n", err)
    }
    fmt.Printf("yaml struct: %+v\n", ty)
}
