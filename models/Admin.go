package models

import (
	"sync"
	"time"
)

type AdminLoginJsonBind struct {
	UserName string `json:"username" binding:"required"`
	PassWord string `json:"password" binding:"required"`
}

type TokenInfo struct {
	Token     string
	ExpiresAt int64
}

type TokenStore struct {
	store sync.Map
}

func (ts *TokenStore) Add(userID string, token TokenInfo) {
	ts.store.Store(userID, token)
}

func (ts *TokenStore) Get(userID string) (TokenInfo, bool) {
	value, ok := ts.store.Load(userID)
	if !ok {
		return TokenInfo{}, false
	}
	tokenInfo, ok := value.(TokenInfo)
	return tokenInfo, ok
}

func (ts *TokenStore) Delete(userID string) {
	ts.store.Delete(userID)
}

func (ts *TokenStore) Clean() {
	ts.store.Range(func(key, value interface{}) bool {
		tokenInfo := value.(TokenInfo)
		if time.Now().After(time.Unix(tokenInfo.ExpiresAt, 0)) {
			ts.store.Delete(key)
		}
		return true
	})
}
