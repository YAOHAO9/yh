[TOC]
### pine
pine 架构理念来源于pomelo(pinus)，以最简洁的方式实现go语言版本的pomelo，希望pine可以站在巨人的肩膀上走的更远
### Getting started

####　初始化项目
```bash
go mod init example
go get github.com/YAOHAO9/pine@0.8
```

#### config.yaml
```
Server:
  SystemName: "awesome_pine"
  ID: connector-0 # 唯一ID
  Kind: connector # 服务器类型
  Host: 127.0.0.1 
  Port: 3014 
  Token: ksYNdrAo # 集群认证Token
  LogType: "Console" # Console/File 
  LogLevel: "Debug" # Debug/Info/Warn/Error

Zookeeper:
  Host: 127.0.0.1
  Port: 2181 
```

#### main.go
```
mkdir connector
cd connector
vi main.go
```
```go
package main

import (
	"github.com/YAOHAO9/pine/application"
	"github.com/YAOHAO9/pine/application/config"
	"github.com/YAOHAO9/pine/rpc/context"
	"github.com/YAOHAO9/pine/rpc/handler"
)

func main() {
	app := application.CreateApp()

  // 作为Connector启动
	app.AsConnector(func(uid string, token string, sessionData map[string]interface{}) error {

		if uid == "" || token == "" {
			return errors.New("Invalid token")
		}
		sessionData[token] = token

		return nil
	})

  // 注册处理客户端请求的Handler
	app.RegisteHandler("handler", func(respCtx *context.RPCCtx) *handler.Resp {
		return &handler.Resp{
			Data: config.GetServerConfig().ID + ": 收到Handler消息",
		}
	})
  // 注册系统间调用的RPC
	app.RegisteRemoter("rpc", func(respCtx *context.RPCCtx) *handler.Resp {

		return &handler.Resp{
			Data: config.GetServerConfig().ID + ": 收到Rpc消息",
		}
  })

	app.Start()
}
```

#### 启动一个connector服务
```
go run connector/main.go #一个最简单的connector已经启动好了
```

#### 连接到Connector
```
ws://127.0.0.1:3014?id=hao&token=xxxxxxx 
```

#### 连接测试
```bash
npm install ts-node typescript -g
ts-node index.ts
```

#### index.ts
```typescript
import * as WebSocket from 'ws'

const requestMap = {}
const ws = new WebSocket('ws://127.0.0.1:3014?id=hao&token=ksYNdrAo')
ws.onopen = async (_: WebSocket.OpenEvent) => {
    console.warn('已连接')

    for (let index = 1; index < 100; index++) {
        sendMessage(index).then(data => {
            console.log(data)
        })
    }
}

function sendMessage(index) {
    return new Promise((resolve,) => {
        ws.send(JSON.stringify({
            Route: 'connector.handler',
            RequestID: index,
            Data: { RequestID: index }
        }))

        requestMap[index] = (data) => {
            resolve(data)
        }
    })
}

ws.onmessage = (event: WebSocket.MessageEvent) => {
    const data = JSON.parse(event.data.toString())
    if (data.RequestID) {
        const cb = requestMap[data.RequestID]
        delete requestMap[data.RequestID]
        cb(data)
    } else {
        console.warn(data)
    }
}

ws.onclose = (event: WebSocket.CloseEvent) => {
    console.warn('连接被关闭', event.reason)
}

ws.onerror = (event: WebSocket.ErrorEvent) => {
    console.error(event.message)
}
```

#### zookeeper docker-compose.yml
(账号：admin 密码：admin)
```
version: '3.1'

services:
  zoo1:
    image: zookeeper
    hostname: zoo1
    ports:
      - 2181:2181
    environment:
      ZOO_MY_ID: 1
      ZOO_SERVERS: server.1=0.0.0.0:2888:3888;2181 server.2=zoo2:2888:3888;2181 server.3=zoo3:2888:3888;2181

  zoo2:
    image: zookeeper
    hostname: zoo2
    ports:
      - 2182:2181
    environment:
      ZOO_MY_ID: 2
      ZOO_SERVERS: server.1=zoo1:2888:3888;2181 server.2=0.0.0.0:2888:3888;2181 server.3=zoo3:2888:3888;2181

  zoo3:
    image: zookeeper
    hostname: zoo3
    ports:
      - 2183:2181
    environment:
      ZOO_MY_ID: 3
      ZOO_SERVERS: server.1=zoo1:2888:3888;2181 server.2=zoo2:2888:3888;2181 server.3=0.0.0.0:2888:3888;2181

  node-zk-browser:
    image: fify/node-zk-browser
    hostname: node-zk-browser
    ports:
      - "3000:3000"
    environment:
      ZK_HOST: zoo1:2181
```

### 更多
##### 添加子游戏
##### 添加路由
##### 添加filter
##### docker-compose
##### k8s
