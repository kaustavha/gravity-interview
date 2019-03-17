package main

import (
	"net/http"
	"time"
)

//Authenticator interface for Authenticator package use
type Authenticator interface {
	IsAuthenticated(r *http.Request) bool
	CleanupExpiredTokens() error
	Upgrade(r *http.Request) ([]byte, error)
	GetInfo(r *http.Request) ([]byte, error)
	LoginAdmin(email string, password string) (string, time.Time, error)
	LogoutAdmin(r *http.Request) (string, time.Time, error)
}
