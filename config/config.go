package config

import (
	"encoding/json"
	"io/ioutil"
)

// ==========================================
// ServerConfig
// ==========================================
var serverCOnfig *ServerConfig

// ServerConfig 服务器配置 配置文件
type ServerConfig struct {
	SystemName string
	Type       string
	Name       string
	Host       string
	Port       string
	Token      string
}

// InitServerConfig 获取服务器配置
func InitServerConfig(path string) *ServerConfig {
	serverCOnfig = &ServerConfig{}
	r, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	json.Unmarshal(r, serverCOnfig)
	return serverCOnfig
}

// GetServerConfig 获取服务器配置
func GetServerConfig() *ServerConfig {
	return serverCOnfig
}

// ==========================================
// ZkConfig
// ==========================================
var zkConfig *ZkConfig

// ZkConfig zk 配置文件
type ZkConfig struct {
	Host string
	Port string
}

// InitZkConfig 初始化
func InitZkConfig(path string) *ZkConfig {
	zkConfig = &ZkConfig{}
	r, e := ioutil.ReadFile(path)
	if e != nil {
		panic(e)
	}
	json.Unmarshal(r, zkConfig)
	return zkConfig
}

// GetZkConfig 获取zk配置
func GetZkConfig() *ZkConfig {
	return zkConfig
}
