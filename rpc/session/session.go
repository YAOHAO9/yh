package session

import "sync"

var mutex sync.Mutex

// Session of connection
type Session struct {
	UID  string // User id
	CID  string // Connector id
	Data map[string]interface{}
}

// Get a value from session
func (session Session) Get(key string) interface{} {
	return session.Data[key]
}

// Set a value to session
func (session Session) Set(key string, v interface{}) {
	mutex.Lock()
	if session.Data == nil {
		session.Data = make(map[string]interface{})
	}
	session.Data[key] = v
	mutex.Unlock()
}
