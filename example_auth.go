package main

import (
	"encoding/json"
	"net/http"

	"github.com/kaustavha/gravity-interview/src/authenticator"
)

// Example of what a login handler will now look like
func ExampleLoginHandler(w http.ResponseWriter, r *http.Request) {
	db := GetDBConn()
	authenticator, err := authenticator.NewAuthenticator(
		"accid",
		"email",
		"pass",
		100,
		[]byte("signingkey"),
		db,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loginHandler(authenticator, w, r)
}

// Example of our login handler helper func
func loginHandler(authenticator *authenticator.Authenticator, w http.ResponseWriter, r *http.Request) {
	// get cookie
	// if active user then update token deets internally and save
	// else try to find in db and set in session
	// else create default and set in session

	c, err := r.Cookie(defaultCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	var creds Credentials

	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// If the structure of the body is wrong, return an HTTP error
		return
	}

	pass, err := authenticator.DecodeAndCheckCreds(creds.Password, creds.Email)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}

	acc := authenticator.Login(sessionToken, pass)

	http.SetCookie(w, &http.Cookie{
		Name:    defaultCookieName,
		Value:   acc.SessionToken,
		Expires: acc.SessionExpiry,
	})

	w.WriteHeader(http.StatusOK)
}
