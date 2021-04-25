package serializer

import (
	"encoding/json"

	"github.com/YAOHAO9/pine/logger"
	"github.com/golang/protobuf/proto"
)

// ToBytes encode anything to []byte
func ToBytes(v interface{}) []byte {
	msesage, ok := v.(proto.Message)
	if ok {
		bytes, err := proto.Marshal(msesage)
		if err != nil {
			logger.Error("Proto消息encode失败", err)
			return []byte{}
		}
		return bytes
	}

	bytes, err := json.Marshal(v)
	if err != nil {
		logger.Error("JSON消息encode失败", err)

		return []byte{123, 125}
	}
	return bytes
}

// FromBytes decode []byte to interface{}
func FromBytes(bytes []byte, v interface{}) error {

	message, ok := v.(proto.Message)
	if ok {
		err := proto.Unmarshal(bytes, message)
		return err
	}

	err := json.Unmarshal(bytes, v)
	return err
}
