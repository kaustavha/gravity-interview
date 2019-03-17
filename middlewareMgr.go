package main

import (
	"fmt"
	"net/http"
	"reflect"
)

//MiddlewareManager class to manage different middlerwares and chaining
type MiddlewareManager struct {
	a Authenticator
}

//GetNewMiddlewareManager returns a new MiddlewareManager
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
		err := m.a.CleanupExpiredTokens()
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, "Error clearing tokens", http.StatusInternalServerError)
		}
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
