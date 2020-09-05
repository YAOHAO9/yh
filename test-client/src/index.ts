
import * as WebSocket from 'ws'

const requestMap = {}
const ws = new WebSocket('ws://127.0.0.1:3110?id=hao&token=ksYNdrAo')
ws.onopen = async (_: WebSocket.OpenEvent) => {
    console.warn('已连接')

    for (let index = 1; index < 100; index++) {
        const data = await sendMessage(index)
        console.log(data)
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
        console.log(data.RequestID)
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