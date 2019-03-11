package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware_withAuthcheckHandler_Success(t *testing.T) {
	InitAuth()
	CreateDBConn()
	req, err := http.NewRequest("GET", "/api/authcheck", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(applyMiddlewares(AuthcheckHandler))
	handler.ServeHTTP(rr, req)

	if err != nil {
		t.Errorf("Did not expect an error response : err: %v ; code: %v", err, rr.Code)
	}
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected: %v but got %v", http.StatusUnauthorized, rr.Code)
	}
}
