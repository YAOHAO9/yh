import Pine from 'pine-client'

(async () => {

    const pine = Pine.init()
    await pine.connect(`ws://127.0.0.1:3014?id=hao&token=ksYNdrAo`)

    pine.on('connector.onMsg', (data) => {
        console.warn('onMsg', data)
    })

    pine.on('connector.onMsgJSON', (data) => {
        console.warn('onMsgJSON', data)
    })

    // const requestData = { Name: 'Proto request', Age: 18 }
    // for (let i = 0; i < 100; i++) {
    //     pine.request('connector.handler', requestData, (response) => {
    //         console.warn('Response:', response)
    //     })
    // }

    await pine.fetchProto('connector')
    const requestDataJSON = { Name: 'JSON request', hahahahah: 18 }
    pine.request('connector.handler', requestDataJSON, (data) => {
        console.warn(data)
    })
    // pine.request('connector.FetchProto__', 'ab23', (data) => {
    //     Object.keys(data).forEach(key => {
    //         data[key] = JSON.parse(data[key])
    //     })

    //     console.warn(JSON.stringify(data, null, 2))
    // })
})()



