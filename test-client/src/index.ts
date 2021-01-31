import Pine from 'pine-client'
import { Middleware } from 'pine-client/lib/common'

(async () => {

    const pine = Pine.init()
    await pine.connect(`ws://127.0.0.1?id=${Math.random()}&token=ksYNdrAo`)


    pine.on('game1.onMsg', (data) => {
        console.warn('game1.onMsg', data)
    })

    const requestDataJSON = { Name: 'JSON request', hahahahah: 18 }


    await pine.fetchCompressMetadata('game1')
    const result1 = await pine.request('game1.handler', requestDataJSON)
    console.log(result1)
    // 中间件1
    const middleware1: Middleware = (data) => {
        if (data.Code === 200) {
            console.warn(data.Message)
            return true
        }
        return false
    }

    // 中间件2
    const middleware2: Middleware = (data) => {
        if (data.Code.toString().startsWith('4')) {
            console.error(data.Message)
            return false
        }
        return true
    }

    // 加入中间件并发
    const result2 = await pine.request('game1.handler', requestDataJSON, middleware1, middleware2)
    console.log(result2)

})()



