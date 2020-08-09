
import * as WebSocket from 'ws'

const ws = new WebSocket('ws://127.0.0.1:3113?id=hao&token=123123')
ws.onopen = (_: WebSocket.OpenEvent) => {
    console.warn('已连接')
    ws.send(JSON.stringify({
        Handler: 'handler',
        Kind: 'ddz',
        Index: 1,
        Data: { a: 1, b: 2, c: 3 }
    }))
}

ws.onmessage = (event: WebSocket.MessageEvent) => {
    console.warn('收到服务端的回复', event.data)
}

ws.onclose = (event: WebSocket.CloseEvent) => {
    console.warn('连接被关闭', event.reason)
}

ws.onerror = (event: WebSocket.ErrorEvent) => {
    console.error(event.message)
}