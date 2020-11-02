package connector

import "errors"

var authFunc = func(uid string, token string) (map[string]interface{}, error) {
	if uid == "" || token == "" {
		return nil, errors.New("认证失败")
	}
	return make(map[string]interface{}), nil
}

// RegisteAuth Registe auth func
func RegisteAuth(auth func(uid string, token string) (map[string]interface{}, error)) {
	authFunc = auth
}
