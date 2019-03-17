package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kaustavha/gravity-interview/src/authenticator"
)

// Testing philosophy: Test the top level to-be-consumed API, and test for the different cases handled by internals
func TestAuthMiddleware_withAuthcheckHandler_Success(t *testing.T) {
	// setup
	db, err := createDBConn()
	a, err := authenticator.NewAuthenticator(
		accountID,
		email,
		hashedPass,
		maxUsers,
		[]byte(signingKey),
		maxUsersUpgraded,
		defaultCookieName,
		db,
	)
	m := GetNewMiddlewareService(a)

	as := GetNewAuthService(a, defaultCookieName)
	// end setup

	req, err := http.NewRequest("GET", "/api/authcheck", nil)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(m.applyMiddlewares(as.AuthcheckHandler))
	handler.ServeHTTP(rr, req)

	if err != nil {
		t.Errorf("Did not expect an error response : err: %v ; code: %v", err, rr.Code)
	}
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected: %v but got %v", http.StatusUnauthorized, rr.Code)
	}
}
