package zookeeper

import (
	"encoding/json"
	"fmt"
	"time"
	"trial/config"

	"github.com/samuel/go-zookeeper/zk"
)

var zkClient *ZkClient

// Start zookeeper
func Start() {

	// 读取配置文件
	zkConfig := config.GetZkConfig()

	// 建立连接
	client, _, err := zk.Connect([]string{zkConfig.Host}, time.Second)
	zkClient = &ZkClient{client: client}
	if err != nil {
		panic(err)
	}

	// 初始化节点
	initNode()

	// 监听节点变化
	watch()
}

// 初始化节点
func initNode() {

	// 服务器配置
	serverConfig := config.GetServerConfig()

	// 检查根节点是否存在，不存在则创建
	rootPath := fmt.Sprint("/", serverConfig.SystemName)
	if !zkClient.exists(rootPath) {
		zkClient.create(rootPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
	}

	// 检查子节点是否存在，不存在则创建
	subPath := fmt.Sprint(rootPath, "/", serverConfig.Kind)
	if !zkClient.exists(subPath) {
		zkClient.create(subPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
	}

	// 解析服务器配置信息
	nodeInfo, err := json.Marshal(serverConfig)
	if err != nil {
		panic(err)
	}

	// 检查服务器数据节点是否存在，不存在则创建
	nodePath := fmt.Sprint(subPath, "/", serverConfig.ID)
	if zkClient.exists(nodePath) {
		zkClient.set(nodePath, nodeInfo, 10)
	} else {
		zkClient.create(nodePath, nodeInfo, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	}
	zkClient.serverID = serverConfig.ID
	fmt.Println("Node created:", nodePath)
}

func watch() {
	// 服务器配置
	serverConfig := config.GetServerConfig()
	path := fmt.Sprint("/", serverConfig.SystemName, "/", serverConfig.Kind)

	for {
		// 遍历所有的serverID
		serverIDs, _, eventChan, err := zkClient.client.ChildrenW(path)
		if err != nil {
			panic(err)
		}
		// 监听每个server的情况
		for _, serverID := range serverIDs {
			if _, ok := watchingServerMap[serverID]; ok {
				// 如果已经建立果监听则跳过
				continue
			}

			go func(serverID string) {
				for {
					// 监听服务器变化
					data, _, eventChan, err := zkClient.client.GetW(fmt.Sprint(path, "/", serverID))
					if err != nil {
						panic(err)
					}
					// 解析服务器信息
					serverConfig := &config.ServerConfig{}
					err = json.Unmarshal(data, serverConfig)
					if err != nil {
						panic(err)
					}
					// 保存服务器信息
					watchingServerMap[serverID] = serverConfig
					// 没有新事件，则阻塞
					<-eventChan
				}
			}(serverID)
		}
		// 没有新事件，则阻塞
		<-eventChan
	}
}
