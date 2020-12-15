package message

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
)

// ToBytes decode anything to []byte
func ToBytes(v interface{}) []byte {
	msesage, ok := v.(proto.Message)
	if ok {
		bytes, err := proto.Marshal(msesage)
		if err != nil {
			logrus.Error("Proto消息encode失败")
			return []byte("Proto消息encode失败")
		}
		return bytes
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		logrus.Error("JSON消息encode失败")

		return []byte("JSON消息encode失败")
	}
	return bytes

}
