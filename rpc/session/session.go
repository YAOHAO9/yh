package session

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
	session.Data[key] = v
}
