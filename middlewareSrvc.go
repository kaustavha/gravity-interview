package main

import (
	"fmt"
	"net/http"
	"reflect"
)

//MiddlewareService class to manage different middlerwares and chaining
type MiddlewareService struct {
	a Authenticator
}

//GetNewMiddlewareService returns a new MiddlewareManager
func GetNewMiddlewareService(a interface{}) *MiddlewareService {
	return &MiddlewareService{
		a: reflect.ValueOf(a).Interface().(Authenticator),
	}
}

func (m *MiddlewareService) getWrappedLoginHandler(LoginHandler http.HandlerFunc) http.HandlerFunc {
	return m.loggingMiddleware(m.cleanupExpiredTokensMiddleware(LoginHandler))
}

func (m *MiddlewareService) getWrappedIOTDataHandler(IOTDataHandler http.HandlerFunc) http.HandlerFunc {
	return m.loggingMiddleware(IOTDataHandler)
}

func (m *MiddlewareService) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if found := m.a.IsAuthenticated(r); found {
			next(w, r)
			return
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func (m *MiddlewareService) cleanupExpiredTokensMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := m.a.CleanupExpiredTokens()
		if err != nil {
			http.Error(w, "Error clearing tokens", http.StatusInternalServerError)
		}
		next(w, r)
		return
	}
}

func (m *MiddlewareService) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("API endpoint hit: ", r.RequestURI)
		next(w, r)
		return
	}
}

func (m *MiddlewareService) applyMiddlewares(next http.HandlerFunc) http.HandlerFunc {
	return m.loggingMiddleware(
		m.cleanupExpiredTokensMiddleware(
			m.authMiddleware(next)))
}
