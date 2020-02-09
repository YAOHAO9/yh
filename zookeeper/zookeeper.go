package zookeeper

import (
	"encoding/json"
	"fmt"
	"time"
	"trial/config"

	"github.com/samuel/go-zookeeper/zk"
)

// ZkClient custom
type ZkClient struct {
	client *zk.Conn
}

var zkClient *ZkClient

func (zkClient ZkClient) exists(path string) bool {
	ok, _, err := zkClient.client.Exists(path)
	if err != nil {
		panic(err)
	}
	return ok
}

func (zkClient ZkClient) create(path string, data []byte, flags int32, acl []zk.ACL) string {
	path, err := zkClient.client.Create(path, data, flags, acl)
	if err != nil {
		panic(err)
	}
	return path
}

func (zkClient ZkClient) set(path string, data []byte, version int32) {
	zkClient.client.Set(path, data, version)
}

// Start zookeeper
func Start() {

	// 读取配置文件
	zkConfig := config.GetZkConfig()

	// 建立连接
	client, _, err := zk.Connect([]string{zkConfig.Host}, time.Second*5)
	zkClient = &ZkClient{client}
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

	// 根节点
	rootPath := fmt.Sprint("/", serverConfig.SystemName)
	if !zkClient.exists(rootPath) {
		zkClient.create(rootPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
	}

	// 子路径
	subPath := fmt.Sprint(rootPath, "/", serverConfig.Type)
	if !zkClient.exists(subPath) {
		zkClient.create(subPath, []byte{}, 0, zk.WorldACL(zk.PermAll))
	}

	// 节点
	nodePath := fmt.Sprint(subPath, "/", serverConfig.Name)
	nodeInfo, err := json.Marshal(serverConfig)
	if err != nil {
		panic(err)
	}
	if zkClient.exists(subPath) {
		zkClient.set(nodePath, nodeInfo, 10)
	} else {
		zkClient.create(nodePath, nodeInfo, zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	}
	fmt.Println("Node created:", nodePath)
}

func watch() {
	// zk.WithEventCallback(func(event zk.Event) {
	// 	fmt.Println("*******************")
	// 	fmt.Println("path:", event.Path)
	// 	fmt.Println("type:", event.Type.String())
	// 	fmt.Println("state:", event.State.String())
	// 	fmt.Println("-------------------")
	// })
	// 服务器配置
	serverConfig := config.GetServerConfig()
	path := fmt.Sprint("/", serverConfig.SystemName, "/", serverConfig.Type)
	_, _, eventChan, err := zkClient.client.ChildrenW(path)
	if err != nil {
		panic(err)
	}

	for {
		event := <-eventChan
		fmt.Println("*******************")
		fmt.Println("path:", event.Path)
		fmt.Println("type:", event.Type.String())
		fmt.Println("state:", event.State.String())
		fmt.Println("-------------------")
	}
}
