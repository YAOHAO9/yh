package connector

import (
	"github.com/sirupsen/logrus"
)

// CompressStore 压缩仓库
var compressStore = make(map[string]map[string]interface{})

// SetCompressData 保存
func SetCompressData(serverKind string, data map[string]interface{}) {
	compressStore[serverKind] = data
}

// GetCompressData 获取
func GetCompressData(serverKind string) map[string]interface{} {
	result, ok := compressStore[serverKind]
	if !ok {
		logrus.Warn("无法获取压缩配置文件 serverKine:", serverKind)
		return make(map[string]interface{})
	}
	return result
}
