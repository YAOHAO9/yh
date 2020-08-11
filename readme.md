## go
``` 
go get  // install modules
sh start.sh // start the app
```

## 启动一个服务
```
eg.
go run main.go -c -p="3110" -k="connector" -i="connector-1" // 启动一个Connector

-c 是否是connector服
-p 端口
-k 服务类型
-i 服务唯一标识
-t Token
-H 服务器IP  默认:127.0.0.1

-zh zookeeper host
-zp zookeeper post

-h help
-s 系统名称
```

## 连接到Connector
ws://127.0.0.1:3110?id=hao&token=xxxxxxx

## zookeeper docker-compose.yml
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