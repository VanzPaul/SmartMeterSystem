/* package utils

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
*/

package utils

import (
	"container/list"
	"sync"
)

type LoginData struct {
	SessionToken, CSRFToken, Role string
}

type Store struct {
	// Web authentication
	WebUsers        map[string]LoginData
	WebSessions     map[string]string
	webUserOrder    *list.List
	webSessionOrder *list.List
	maxWebUsers     int
	maxWebSessions  int
	webUserMutex    sync.Mutex
	webSessionMutex sync.Mutex

	// Meter authentication
	MeterUsers        map[string]LoginData
	MeterSessions     map[string]string
	meterUserOrder    *list.List
	meterSessionOrder *list.List
	maxMeterUsers     int
	maxMeterSessions  int
	meterUserMutex    sync.Mutex
	meterSessionMutex sync.Mutex
}

var globalStore *Store

// GetStore initializes the store with separate limits for web and meter authentication
func GetStore(maxWebUsers, maxWebSessions, maxMeterUsers, maxMeterSessions int) *Store {
	if globalStore == nil {
		globalStore = &Store{
			WebUsers:        make(map[string]LoginData),
			WebSessions:     make(map[string]string),
			webUserOrder:    list.New(),
			webSessionOrder: list.New(),
			maxWebUsers:     maxWebUsers,
			maxWebSessions:  maxWebSessions,

			MeterUsers:        make(map[string]LoginData),
			MeterSessions:     make(map[string]string),
			meterUserOrder:    list.New(),
			meterSessionOrder: list.New(),
			maxMeterUsers:     maxMeterUsers,
			maxMeterSessions:  maxMeterSessions,
		}
	}
	return globalStore
}

// Web user methods
func (s *Store) AddWebUser(username string, data LoginData) {
	s.webUserMutex.Lock()
	defer s.webUserMutex.Unlock()

	if len(s.WebUsers) >= s.maxWebUsers {
		oldest := s.webUserOrder.Front()
		if oldest != nil {
			delete(s.WebUsers, oldest.Value.(string))
			s.webUserOrder.Remove(oldest)
		}
	}

	s.WebUsers[username] = data
	s.webUserOrder.PushBack(username)
}

func (s *Store) GetWebUser(username string) (LoginData, bool) {
	s.webUserMutex.Lock()
	defer s.webUserMutex.Unlock()

	data, exists := s.WebUsers[username]
	return data, exists
}

// Web session methods
func (s *Store) AddWebSession(sessionID, username string) {
	s.webSessionMutex.Lock()
	defer s.webSessionMutex.Unlock()

	if len(s.WebSessions) >= s.maxWebSessions {
		oldest := s.webSessionOrder.Front()
		if oldest != nil {
			delete(s.WebSessions, oldest.Value.(string))
			s.webSessionOrder.Remove(oldest)
		}
	}

	s.WebSessions[sessionID] = username
	s.webSessionOrder.PushBack(sessionID)
}

func (s *Store) GetWebSession(sessionID string) (string, bool) {
	s.webSessionMutex.Lock()
	defer s.webSessionMutex.Unlock()

	username, exists := s.WebSessions[sessionID]
	return username, exists
}

// Meter user methods
func (s *Store) AddMeterUser(username string, data LoginData) {
	s.meterUserMutex.Lock()
	defer s.meterUserMutex.Unlock()

	if len(s.MeterUsers) >= s.maxMeterUsers {
		oldest := s.meterUserOrder.Front()
		if oldest != nil {
			delete(s.MeterUsers, oldest.Value.(string))
			s.meterUserOrder.Remove(oldest)
		}
	}

	s.MeterUsers[username] = data
	s.meterUserOrder.PushBack(username)
}

func (s *Store) GetMeterUser(username string) (LoginData, bool) {
	s.meterUserMutex.Lock()
	defer s.meterUserMutex.Unlock()

	data, exists := s.MeterUsers[username]
	return data, exists
}

// Meter session methods
func (s *Store) AddMeterSession(sessionID, username string) {
	s.meterSessionMutex.Lock()
	defer s.meterSessionMutex.Unlock()

	if len(s.MeterSessions) >= s.maxMeterSessions {
		oldest := s.meterSessionOrder.Front()
		if oldest != nil {
			delete(s.MeterSessions, oldest.Value.(string))
			s.meterSessionOrder.Remove(oldest)
		}
	}

	s.MeterSessions[sessionID] = username
	s.meterSessionOrder.PushBack(sessionID)
}

func (s *Store) GetMeterSession(sessionID string) (string, bool) {
	s.meterSessionMutex.Lock()
	defer s.meterSessionMutex.Unlock()

	username, exists := s.MeterSessions[sessionID]
	return username, exists
}

// In models/store.go
func (s *Store) DeleteWebSession(sessionID string) {
	s.webSessionMutex.Lock()
	defer s.webSessionMutex.Unlock()
	delete(s.WebSessions, sessionID)
}
