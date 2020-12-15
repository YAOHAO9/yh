package session

import "sync"

var mutex sync.Mutex

// Get a value from session
func (session *Session) Get(key string) interface{} {
	return session.Data[key]
}

// Set a value to session
func (session *Session) Set(key string, v string) {
	mutex.Lock()
	if session.Data == nil {
		session.Data = make(map[string]string)
	}
	session.Data[key] = v
	mutex.Unlock()
}
