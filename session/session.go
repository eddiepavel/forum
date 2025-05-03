package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// Session represents a user session.
type Session struct {
	ID        string
	Data      map[string]interface{} // Used for flash messages and other session data.
	ExpiresAt time.Time
}

// SessionStore manages all sessions in memory.
type SessionStore struct {
	sessions        map[string]*Session
	mutex           sync.RWMutex
	SessionDuration time.Duration
}

// NewSessionStore creates a new session store and starts a cleanup ticker.
// duration: how long a session is valid.
// cleanupInterval: how often to check for expired sessions.
func NewSessionStore(duration, cleanupInterval time.Duration) *SessionStore {
	store := &SessionStore{
		sessions:        make(map[string]*Session),
		SessionDuration: duration,
	}
	go store.cleanupTicker(cleanupInterval)
	return store
}

// cleanupTicker periodically calls cleanup() to remove expired sessions.
func (s *SessionStore) cleanupTicker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		s.cleanup()
	}
}

// cleanup removes expired sessions from the store.
func (s *SessionStore) cleanup() {
	fmt.Println("started cleaning session")
	s.mutex.Lock()
	defer s.mutex.Unlock()
	now := time.Now()
	for id, sess := range s.sessions {
		if sess.ExpiresAt.Before(now) {
			delete(s.sessions, id)
		}
	}
}

// CreateSession generates a new session, sets its expiration, and adds it to the store.
func (s *SessionStore) CreateSession() *Session {
	id := generateSessionID()
	session := &Session{
		ID:        id,
		Data:      make(map[string]interface{}),
		ExpiresAt: time.Now().Add(s.SessionDuration),
	}
	s.mutex.Lock()
	s.sessions[id] = session
	s.mutex.Unlock()
	return session
}

// GetSession retrieves a session by its ID.
// It also checks if the session is expired and, if so, removes it.
func (s *SessionStore) GetSession(id string) (*Session, bool) {
	s.mutex.RLock()
	session, exists := s.sessions[id]
	s.mutex.RUnlock()
	if !exists || session.ExpiresAt.Before(time.Now()) {
		// If expired or not found, remove it if necessary.
		s.mutex.Lock()
		delete(s.sessions, id)
		s.mutex.Unlock()
		return nil, false
	}
	return session, true
}

// RemoveSession manually removes a session from the store.
func (s *SessionStore) RemoveSession(id string) {
	s.mutex.Lock()
	delete(s.sessions, id)
	s.mutex.Unlock()
}

// RefreshSession extends a session's expiration time.
func (s *SessionStore) RefreshSession(id string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if session, exists := s.sessions[id]; exists {
		session.ExpiresAt = time.Now().Add(s.SessionDuration)
		return true
	}
	return false
}

// generateSessionID creates a random 32-character hexadecimal session ID.
func generateSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback: use the current timestamp if random generation fails.
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

// Flash message utilities on the Session type:

// SetFlash stores a flash message under the given key.
func (sess *Session) SetFlash(key string, value interface{}) {
	if sess.Data == nil {
		sess.Data = make(map[string]interface{})
	}
	sess.Data[key] = value
}

// GetFlash retrieves and removes a flash message from the session.
func (sess *Session) GetFlash(key string) (interface{}, bool) {
	if sess.Data == nil {
		return nil, false
	}
	val, exists := sess.Data[key]
	if exists {
		// Remove after reading
		delete(sess.Data, key)
	}
	return val, exists
}
