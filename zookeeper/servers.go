package zookeeper

import (
	"math/rand"
	"trial/config"
)

// 是否正在监听
var watchingServerMap = make(map[string]*config.ServerConfig)

// GetServersByType 获取统一种类的所有server
func GetServersByType(kind string) []*config.ServerConfig {
	servers := make([]*config.ServerConfig, 0)

	for _, server := range servers {
		if server.Kind == kind {
			servers = append(servers, server)
		}
	}
	return servers
}

// GetServerByID 根据ID获取server
func GetServerByID(id string) *config.ServerConfig {
	servers := make([]*config.ServerConfig, 0)

	for _, server := range servers {
		if server.ID == id {
			return server
		}
	}
	return nil
}

// GetRandServerByType 根据总累随机获取一个server
func GetRandServerByType(kind string) *config.ServerConfig {
	servers := make([]*config.ServerConfig, 0)

	for _, server := range servers {
		if server.Kind == kind {
			servers = append(servers, server)
		}
	}

	return servers[rand.Intn(len(servers))]
}
