package application

import (
	"github.com/YAOHAO9/pine/connector"
	"github.com/spf13/viper"
)

// AsConnector 作为Connector 启动
func (app Application) AsConnector(authFunc func(uid string, token string, sessionData map[string]interface{}) error) {
	viper.Set("Server.IsConnector", true)
	connector.RegisteAuth(authFunc)
}
