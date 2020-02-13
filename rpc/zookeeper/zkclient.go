package zookeeper

import (
	"github.com/samuel/go-zookeeper/zk"
)

// ZkClient custom
type ZkClient struct {
	client   *zk.Conn
	serverID string
}

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

// Close zk client
func (zkClient ZkClient) Close() {
	zkClient.Close()
}
