package zookeeper

import (
	"github.com/samuel/go-zookeeper/zk"
)

// ZkClient custom
type ZkClient struct {
	conn     *zk.Conn
	serverID string
}

func (client ZkClient) exists(path string) bool {
	ok, _, err := client.conn.Exists(path)
	if err != nil {
		panic(err)
	}
	return ok
}

func (client ZkClient) create(path string, data []byte, flags int32, acl []zk.ACL) string {
	path, err := client.conn.Create(path, data, flags, acl)
	if err != nil {
		panic(err)
	}
	return path
}

func (client ZkClient) set(path string, data []byte, version int32) {
	client.conn.Set(path, data, version)
}

// Close zk client
func (client ZkClient) Close() {
	client.Close()
}
