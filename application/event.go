package application

import "github.com/YAOHAO9/pine/application/channelservice/eventcompress"

// // AddEventCompressRecord 添加路由压缩记录
// func (app *Application) AddEventCompressRecord(eventName string) {
// 	eventcompress.AddEventRecord(eventName)
// }

// AddEventCompressRecords 添加路由压缩记录
func (app *Application) AddEventCompressRecords(eventNames ...string) {
	for _, eventName := range eventNames {
		eventcompress.AddEventRecord(eventName)
	}
}
