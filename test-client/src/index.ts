
import * as WebSocket from 'ws'

const ws = new WebSocket('ws://127.0.0.1:3110?id=hao&token=ksYNdrAo')
ws.onopen = (_: WebSocket.OpenEvent) => {
    console.warn('已连接')
    Array.from({ length: 100 }).forEach((_, index) => {
        ws.send(JSON.stringify({
            Route: 'connector.handler',
            RequestID: index,
            Data: { a: 1, b: 2, c: 3 }
        }))
    })
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