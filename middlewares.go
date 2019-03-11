package main

import (
	"fmt"
	"net/http"
)

func getWrappedIOTDataHandler() http.HandlerFunc {
	return loggingMiddleware(IOTDataHandler)
}

func getWrappedLoginHandler() http.HandlerFunc {
	return loggingMiddleware(cleanupExpiredTokensMiddleware(LoginHandler))
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isAuthenticated(r) {
			next(w, r)
			return
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("API endpoint hit: ", r.RequestURI)
		next(w, r)
		return
	}
}

func cleanupExpiredTokensMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cleanupExpiredTokens(expected)
		next(w, r)
		return
	}
}

func applyMiddlewares(next http.HandlerFunc) http.HandlerFunc {
	return loggingMiddleware(
		cleanupExpiredTokensMiddleware(
			authMiddleware(next)))
}
