### 启动connector服务
```bash
# 启动 connector
export server_port=3014; go run connector/main.go 
# 启动 game1
export server_port=3015; export server_id=game1-0; export server_kind=game1; go run game1/main.go
```