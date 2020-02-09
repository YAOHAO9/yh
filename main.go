package main

import (
	"fmt"
	"net/http"
	WsServer "trial/ws/server"
)

func main() {

	// Handler
	http.HandleFunc("/ws", WsServer.WebSocketHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("你好世界！！！"))
	})

	// ListenAndServe
	port := "8080"
	fmt.Println("Server started http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil)
	fmt.Println(err.Error())
}
