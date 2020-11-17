
import Pine from 'pine-client'

(async () => {
    const pine = Pine.init()
    await pine.connect(`ws://127.0.0.1:3014?id=${process.argv[2]}&token=ksYNdrAo`)
    pine.on('test', (data) => {
        console.warn('Event test:', data)
    })
    pine.request('connector.haha', 111, (data) => {
        console.warn('Response:', data)
    })
})()