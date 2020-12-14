package serverdict

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// // ServerRouter handler map Reocrd
// type ServerRouter struct {
// 	Code              int
// 	HandlerMap        map[string]int
// 	ReverseHandlerMap map[int]string
// 	NotifyMap         map[string]int
// 	ReverseNotifyMap  map[int]string
// }

// // ServerDict ServerDict
// type ServerDict struct {
// 	serverMap        map[string]*ServerRouter
// 	reverseServerMap map[int]string
// }

// func (routerDict *ServerDict) addRecord(serverKind string, handlers []string, notifys []string) {

// 	if _, exist := routerDict.serverMap[serverKind]; exist {
// 		return
// 	}

// 	code := len(routerDict.reverseServerMap) + 1

// 	routerDict.reverseServerMap[code] = serverKind

// 	serverRouter := &ServerRouter{
// 		Code:              code,
// 		HandlerMap:        make(map[string]int, 10),
// 		ReverseHandlerMap: make(map[int]string, 10),
// 		NotifyMap:         make(map[string]int, 10),
// 		ReverseNotifyMap:  make(map[int]string, 10),
// 	}

// 	routerDict.serverMap[serverKind] = serverRouter

// 	for index, handler := range handlers {
// 		serverRouter.HandlerMap[handler] = index + 1
// 		serverRouter.ReverseHandlerMap[index+1] = handler
// 	}

// 	for index, notify := range notifys {
// 		serverRouter.NotifyMap[notify] = index + 1
// 		serverRouter.ReverseNotifyMap[index+1] = notify
// 	}
// }

var kindToCode = make(map[string]byte)
var codeToKind = make(map[byte]string)

// AddRecord add serverKind and serverCode recore
func AddRecord(serverKind string) {
	if _, exist := kindToCode[serverKind]; exist {
		return
	}

	code := byte(len(kindToCode) + 1)
	kindToCode[serverKind] = code
	codeToKind[code] = serverKind
}

// GetKindByCode get serverKind by serverCode
func GetKindByCode(code byte) string {
	if value, exist := codeToKind[code]; exist {
		return value
	}
	return ""
}

// GetCodeByKind get serverCode by serverKind
func GetCodeByKind(serverKind string) byte {
	if value, exist := kindToCode[serverKind]; exist {
		return value
	}
	return 0
}

// ToBytes get json bytes
func ToBytes() []byte {
	bytes, err := json.Marshal(map[string]interface{}{
		"kindToCode": kindToCode,
		"codeToKind": codeToKind,
	})
	if err != nil {
		logrus.Error(err)
		return []byte{123, 125}
	}
	return bytes
}
