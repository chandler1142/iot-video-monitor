package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	ServerIp     string `json:"serverIp"`
	ServerPort   int `json:"serverPort"`
	SysAddrCode  string `json:"sysAddrCode"`
	ServerDomain string `json: serverDomain`
	LocalIp      string `json: localIp`
	LocalPort    int `json: localPort`
	TemplatePath string `json:templatePath`
}

//初始化配置文件对象
func NewConfig(path string) (*Config, error) {
	c := Config{}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Read config file fail: %v \n", err)
		return nil, err
	}
	json.Unmarshal(buf, &c)
	return &c, nil
}
