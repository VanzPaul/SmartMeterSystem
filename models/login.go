package models

// FIXME: remove HashedPassword, AccountNo,
type LoginData struct {
	HashedPassword, AccountNo, SessionToken, CSRFToken string
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
