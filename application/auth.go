package application

import (
	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/connector"
)

// AsConnector 作为Connector 启动
func (app Application) AsConnector(authFunc func(uid string, token string, sessionData map[string]string) error) {
	config.GetServerConfig().IsConnector = true
	connector.RegisteAuth(authFunc)
}
