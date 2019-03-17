package main

import "net/http"

//Authenticator interface for Authenticator package use
type Authenticator interface {
	IsAuthenticated(r *http.Request) bool
	CleanupExpiredTokens() error
	Upgrade(r *http.Request) ([]byte, error)
	GetInfo(r *http.Request) ([]byte, error)
}
