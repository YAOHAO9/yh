import Pine from 'pine-client'
import * as util from 'util'

(async () => {

    const pine = Pine.init()
    await pine.connect(`ws://127.0.0.1:3014?id=hao&token=ksYNdrAo`)

    pine.on('onMsg', (data) => {
        console.warn('onMsg', data)
    })

    const requestData = { name: 'asdf', age: 18 }
    for (let i = 0; i < 10; i++) {
        await new Promise(resolve => {
            setTimeout(() => {
                resolve()
            }, 500);
        })
        pine.request('connector.handler', new util.TextEncoder().encode(JSON.stringify(requestData)), (response) => {
            console.warn('Response:', response)
        })
    }
})()



