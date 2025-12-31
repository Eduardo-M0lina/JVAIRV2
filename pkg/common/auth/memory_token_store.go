package auth

import (
	"context"
	"sync"
	"time"
)

// MemoryTokenStore implementa TokenStore usando un mapa en memoria
type MemoryTokenStore struct {
	tokens map[string]tokenInfo
	mu     sync.RWMutex
}

type tokenInfo struct {
	userID     int64
	expiration time.Time
}

// NewMemoryTokenStore crea una nueva instancia de MemoryTokenStore
func NewMemoryTokenStore() *MemoryTokenStore {
	return &MemoryTokenStore{
		tokens: make(map[string]tokenInfo),
	}
}

// StoreToken almacena un token en memoria
func (s *MemoryTokenStore) StoreToken(ctx context.Context, userID int64, tokenID string, expiration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens[tokenID] = tokenInfo{
		userID:     userID,
		expiration: time.Now().Add(expiration),
	}

	return nil
}

// DeleteToken elimina un token de memoria
func (s *MemoryTokenStore) DeleteToken(ctx context.Context, tokenID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tokens, tokenID)

	return nil
}

// CheckToken verifica si un token existe y no ha expirado
func (s *MemoryTokenStore) CheckToken(ctx context.Context, tokenID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	info, exists := s.tokens[tokenID]
	if !exists {
		return false, nil
	}

	// Verificar si el token ha expirado
	if time.Now().After(info.expiration) {
		// Eliminar token expirado
		go func() {
			s.mu.Lock()
			defer s.mu.Unlock()
			delete(s.tokens, tokenID)
		}()
		return false, nil
	}

	return true, nil
}

// CleanupExpired elimina tokens expirados
func (s *MemoryTokenStore) CleanupExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for tokenID, info := range s.tokens {
		if now.After(info.expiration) {
			delete(s.tokens, tokenID)
		}
	}
}
