package zookeeper

import (
	"github.com/YAOHAO9/pine/logger"
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
		logger.Error(err)
	}
	return ok
}

func (client ZkClient) create(path string, bytes []byte, flags int32, acl []zk.ACL) string {
	path, err := client.conn.Create(path, bytes, flags, acl)
	if err != nil {
		logger.Error(err)
	}
	return path
}

func (client ZkClient) set(path string, bytes []byte, version int32) {
	client.conn.Set(path, bytes, version)
}

// Close zk client
func (client ZkClient) Close() {
	client.Close()
}
