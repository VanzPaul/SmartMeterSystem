package models

type LoginData struct {
	SessionToken, CSRFToken, Role string
}

// Store encapsulates users and sessions data.
type Store struct {
	Users    map[string]LoginData
	Sessions map[string]string
}

// Singleton instance of the Store
var globalStore *Store

// GetStore returns the singleton instance of the Store.
func GetStore() *Store {
	if globalStore == nil {
		globalStore = &Store{
			Users:    make(map[string]LoginData),
			Sessions: make(map[string]string),
		}
	}
	return globalStore
}

// TODO: Add a functionality to store mater users and meter sessions
