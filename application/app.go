package application

import (
	"math/rand"
	"time"

	"github.com/YAOHAO9/yh/application/config"
	"github.com/YAOHAO9/yh/logger"
	RpcServer "github.com/YAOHAO9/yh/rpc/server"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Application app
type Application struct {
}

// Start start application
func (app Application) Start() {
	RpcServer.Start()
}

var app *Application

// CreateApp 创建app
func CreateApp() *Application {

	parseFlag()

	logger.SetLogMode(config.GetServerConfig().LogType)

	if app != nil {
		return app
	}
	app = &Application{}

	return app
}
