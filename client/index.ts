
import * as WebSocket from "ws"

let ws = new WebSocket("ws://127.0.0.1:3113?id=hao&token=123123")
ws.onopen = () => {
    console.warn("已连接")
}

ws.onmessage = (data) => {
    console.warn("收到服务端的回复", data)
}