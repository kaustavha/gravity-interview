package main

import "net/http"

type Authenticator interface {
	IsAuthenticated(r *http.Request) bool
	CleanupExpiredTokens()
	Upgrade(r *http.Request) ([]byte, error)
	GetInfo(r *http.Request) ([]byte, error)
}
