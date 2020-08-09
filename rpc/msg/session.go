package msg

// Session of connection
type Session struct {
	UID  string // User id
	CID  string // Connector id
	Data map[string]interface{}
}

// Get a value from session
func (s Session) Get(key string) interface{} {
	return s.Data[key]
}

// Set a value to session
func (s Session) Set(key string, v interface{}) {
	s.Data[key] = v
}
