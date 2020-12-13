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

// ServerDict int
type serverDict struct {
	KindToCode map[string]byte `json:"kindToCode"`
	CodeToKind map[byte]string `json:"codeToKind"`
}

// AddRecord add serverKind and serverCode recore
func (dict *serverDict) AddRecord(serverKind string) {
	if _, exist := dict.KindToCode[serverKind]; exist {
		return
	}

	code := byte(len(dict.KindToCode) + 1)
	dict.KindToCode[serverKind] = code
	dict.CodeToKind[code] = serverKind
}

// GetKindByCode get serverKind by serverCode
func (dict *serverDict) GetKindByCode(code byte) string {
	if value, exist := dict.CodeToKind[code]; exist {
		return value
	}
	return ""
}

// GetCodeByKind get serverCode by serverKind
func (dict *serverDict) GetCodeByKind(serverKind string) byte {
	if value, exist := dict.KindToCode[serverKind]; exist {
		return value
	}
	return 0
}

// ToBytes get json bytes
func (dict *serverDict) ToBytes() []byte {
	bytes, err := json.Marshal(dict)
	if err != nil {
		logrus.Error(err)
		return []byte{123, 125}
	}
	return bytes
}

// Store serverdict store
var Store = &serverDict{
	KindToCode: make(map[string]byte, 10),
	CodeToKind: make(map[byte]string, 10),
}
