## go
``` 
go get  // install modules
sh start.sh // start the app
```
## 启动一个服务
```
-c 是否是connector服
-p 端口
-k 服务类型
-i 服务唯一标识

go run main.go -c -p="3110" -k="connector" -i="connector-1" // 启动一个Connector
```
## zookeeper