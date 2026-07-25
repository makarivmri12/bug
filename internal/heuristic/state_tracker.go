package heuristic

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

// SessionState represents the state of a user session
type SessionState struct {
	UserID   string
	Tokens   map[string]string
	Cookies  map[string]string
	Headers  map[string]string
	Metadata map[string]interface{}
	mu       sync.RWMutex
}

// StateTracker maintains and tracks session state changes
type StateTracker struct {
	sessions map[string]*SessionState
	mu       sync.RWMutex
	logger   *zap.Logger
}

// NewStateTracker creates a new state tracker
func NewStateTracker(logger *zap.Logger) *StateTracker {
	return &StateTracker{
		sessions: make(map[string]*SessionState),
		logger:   logger,
	}
}

// CreateSession creates a new session
func (st *StateTracker) CreateSession(ctx context.Context, userID string) *SessionState {
	st.mu.Lock()
	defer st.mu.Unlock()

	session := &SessionState{
		UserID:   userID,
		Tokens:   make(map[string]string),
		Cookies:  make(map[string]string),
		Headers:  make(map[string]string),
		Metadata: make(map[string]interface{}),
	}

	st.sessions[userID] = session
	st.logger.Info("session created", zap.String("user_id", userID))

	return session
}

// GetSession retrieves a session
func (st *StateTracker) GetSession(ctx context.Context, userID string) *SessionState {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.sessions[userID]
}

// UpdateSession updates session state
func (st *StateTracker) UpdateSession(ctx context.Context, userID string, tokens map[string]string, cookies map[string]string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	session, ok := st.sessions[userID]
	if !ok {
		return fmt.Errorf("session not found for user: %s", userID)
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	for k, v := range tokens {
		session.Tokens[k] = v
	}
	for k, v := range cookies {
		session.Cookies[k] = v
	}

	st.logger.Info("session updated", zap.String("user_id", userID))
	return nil
}

// DeleteSession removes a session
func (st *StateTracker) DeleteSession(ctx context.Context, userID string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	delete(st.sessions, userID)
	st.logger.Info("session deleted", zap.String("user_id", userID))
	return nil
}
