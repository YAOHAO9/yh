import 'regenerator-runtime/runtime'
import Pine from 'pine-client'

(async () => {

    const pine = Pine.init()
    await pine.connect(`ws://127.0.0.1?id=${Math.random()}&token=ksYNdrAo`)

    pine.on('connector.onMsg', (data) => {
        // console.warn('connector.onMsg', data)
    })

    pine.on('game1.onMsg', (data) => {
        // console.warn('game1.onMsg', data)
    })

    const timesI = 100
    const timesJ = 10000

    await pine.fetchProto('connector') // 第一次访问前先获取protobuf描述文件
    await pine.fetchProto('game1')
    const requestDataJSON = { Name: 'JSON request', hahahahah: 18 }
    for (let i = 0; i < timesI; i++) {
        const tasks = []
        for (let j = 0; j < timesJ; j++) {
            const task = pine.request('game1.handler', requestDataJSON)
            tasks.push(task)
        }
        const datas = await Promise.all(tasks)
        console.warn(await pine.request('game1.handler', requestDataJSON))
        console.warn(await pine.request('connector.handler', requestDataJSON))
    }
    console.warn('finished')
})()



