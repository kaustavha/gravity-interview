package main

import (
	"fmt"
	"net/http"
	"reflect"
)

type MiddlewareManager struct {
	a Authenticator
}

func GetNewMiddlewareManager(a interface{}) *MiddlewareManager {
	return &MiddlewareManager{
		a: reflect.ValueOf(a).Interface().(Authenticator),
	}
}

func (m *MiddlewareManager) getWrappedLoginHandler(LoginHandler http.HandlerFunc) http.HandlerFunc {
	return m.loggingMiddleware(m.cleanupExpiredTokensMiddleware(LoginHandler))
}

func (m *MiddlewareManager) getWrappedIOTDataHandler(IOTDataHandler http.HandlerFunc) http.HandlerFunc {
	return m.loggingMiddleware(IOTDataHandler)
}

func (m *MiddlewareManager) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if found := m.a.IsAuthenticated(r); found {
			next(w, r)
			return
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func (m *MiddlewareManager) cleanupExpiredTokensMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.a.CleanupExpiredTokens()
		next(w, r)
		return
	}
}

func (m *MiddlewareManager) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("API endpoint hit: ", r.RequestURI)
		next(w, r)
		return
	}
}

func (m *MiddlewareManager) applyMiddlewares(next http.HandlerFunc) http.HandlerFunc {
	return m.loggingMiddleware(
		m.cleanupExpiredTokensMiddleware(
			m.authMiddleware(next)))
}
