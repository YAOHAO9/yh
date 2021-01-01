import Pine from 'pine-client'

(async () => {

    const pine = Pine.init()
    await pine.connect(`ws://127.0.0.1:3014?id=${Math.random()}&token=ksYNdrAo`)

    pine.on('connector.onMsg', (data) => {
        console.warn('onMsg', data)
    })

    await pine.fetchProto('connector') // 第一次访问前先获取protobuf描述文件
    const requestDataJSON = { Name: 'JSON request', hahahahah: 18 }
    pine.request('connector.handler', requestDataJSON, (data) => {
        console.warn(data)
    })
})()



