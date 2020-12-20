import Pine from 'pine-client'

(async () => {

    const pine = Pine.init()
    await pine.connect(`ws://127.0.0.1:3014?id=${Math.random()}&token=ksYNdrAo`)

    // pine.on('connector.onMsg1', (data) => {
    //     console.warn('onMsg1', data)
    // })
    // pine.on('connector.onMsg2', (data) => {
    //     console.warn('onMsg2', data)
    // })
    // pine.on('connector.onMsg3', (data) => {
    //     console.warn('onMsg3', data)
    // })
    // pine.on('connector.onMsg4', (data) => {
    //     console.warn('onMsg4', data)
    // })

    // pine.on('connector.onMsgJSON', (data) => {
    //     console.warn('onMsgJSON', data)
    // })

    // pine.on('connector.__Kick__', (data) => {
    //     console.warn('我被踢下线了啊', data)
    // })


    await pine.fetchProto('connector')
    const requestData = { Name: 'Proto request', Age: 18 }
    for (let i = 0; i < 1000; i++) {
        pine.request('connector.handler', requestData, (response) => {
            // console.warn('Response:', response)
        })
    }


    const requestDataJSON = { Name: 'JSON request', hahahahah: 18 }
    pine.request('connector.handler', requestDataJSON, (data) => {
        console.warn(data)
    })
})()



