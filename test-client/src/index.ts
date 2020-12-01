import Pine from 'pine-client'

(async () => {

    const pine = Pine.init()
    await pine.connect(`ws://127.0.0.1:3014?id=hao&token=ksYNdrAo`)

    pine.on('onMsg', (data) => {
        console.warn('onMsg', data)
    })

    const requestData = { name: 'asdf', age: 18 }
    for (let i = 0; i < 2; i++) {
        pine.request('connector.handler', requestData, (response) => {
            console.warn('Response:', response)
        })
    }
})()