package connector

import (
	"github.com/YAOHAO9/pine/logger"
)

var authFunc = func(uid string, token string, sessionData map[string]string) error {
	if uid == "" || token == "" {
		return logger.NewError("认证失败")
	}
	return nil
}

// RegisteAuth Registe auth func
func RegisteAuth(auth func(uid string, token string, sessionData map[string]string) error) {
	authFunc = auth
}
