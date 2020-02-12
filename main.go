package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
	"trial/config"
	WsServer "trial/ws/server"
)

func main() {

	rand.Seed(time.Now().UnixNano())
	serverConfig := config.InitServerConfig("./config/server.json")

	// Handler
	http.HandleFunc("/ws", WsServer.WebSocketHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("你好世界！！！"))
	})

	// ListenAndServe
	fmt.Println("Server started " + serverConfig.Host + ":" + serverConfig.Port)
	err := http.ListenAndServe(":"+serverConfig.Port, nil)
	fmt.Println(err.Error())

}
