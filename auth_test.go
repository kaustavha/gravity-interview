package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaustavha/gravity-interview/src/authenticator"
)

func TestAuthMiddleware_withAuthcheckHandler_Success(t *testing.T) {
	// setup
	db, err := createDBConn()
	a, err := authenticator.NewAuthenticator(
		AccountID,
		Email,
		HashedPass,
		maxUsers,
		[]byte(SigningKey),
		maxUsersUpgraded,
		defaultCookieName,
		db,
	)
	m := GetNewMiddlewareManager(a)
	// end setup

	req, err := http.NewRequest("GET", "/api/authcheck", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(m.applyMiddlewares(AuthcheckHandler))
	handler.ServeHTTP(rr, req)

	if err != nil {
		t.Errorf("Did not expect an error response : err: %v ; code: %v", err, rr.Code)
	}
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected: %v but got %v", http.StatusUnauthorized, rr.Code)
	}
}
