package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var expected Credentials
var mySigningKey []byte

// map of email -> Adminaccs
var activeAccounts map[string]AdminAccount

// map session -> email for lookup in activeAccounts map
var sessionMap map[string]string

// set of active session token strings
var sessionTokens []string

func InitAuth() {
	sessionTokens = []string{}
	activeAccounts = make(map[string]AdminAccount)
	sessionMap = make(map[string]string)
	mySigningKey = []byte(defaultHashedPass)

	expected = Credentials{
		Email:    defaultEmail,
		Password: defaultHashedPass,
	}
}

func AuthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	acc, found := findUserAccountFromActiveToken(r)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	acc.CleanToken()
	w.WriteHeader(http.StatusOK)
}

// StatusOK if already signed in, or after signin based on incoming cookie or email+password
// otherwise frontend needs to prompt for usrname and pass and try again
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		w.WriteHeader(http.StatusOK)
		return
	}
	creds, status := decodeAndCheckCreds(r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		return
	}

	expiry := time.Now().Add(120 * time.Second)
	tokenString := getJWT(creds, expiry, mySigningKey)

	// First check active sessions, then search db
	account, found := activeAccounts[creds.Email]
	if found {
		account.SessionToken = tokenString
		account.SessionExpiry = expiry
	} else {
		found, acc := db.findAdmin(defaultAccountID)
		if found {
			account = *acc
			account.SessionToken = tokenString
			account.SessionExpiry = expiry
		} else {
			account = AdminAccount{
				Email:         creds.Email,
				IsUpgraded:    false,
				SessionExpiry: expiry,
				SessionToken:  tokenString,
				AccountId:     defaultAccountID,
				Users:         0,
				MaxUsers:      maxUsers,
			}
			account.Users = acc.CountAssociatedUsers()
		}
	}

	account.UpdateSelf()
	account.SaveInDB()

	http.SetCookie(w, &http.Cookie{
		Name:    defaultCookieName,
		Value:   account.SessionToken,
		Expires: account.SessionExpiry,
	})

	w.WriteHeader(http.StatusOK)
}

// Helpers

func decodeAndCheckCreds(r *http.Request) (Credentials, int) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		return creds, http.StatusBadRequest
	}

	if CheckPasswordHash(creds.Password, expected.Password) == false {
		return creds, http.StatusUnauthorized
	}

	creds.Password, _ = HashPassword(creds.Password)

	if expected.Email != creds.Email {
		return creds, http.StatusUnauthorized
	}
	return creds, http.StatusOK
}

func isAuthenticated(r *http.Request) bool {
	c, err := r.Cookie(defaultCookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return false
		}
		return false
	}
	sessionToken := c.Value
	found := index(sessionTokens, sessionToken)
	if found != -1 {
		return true
	}
	return false
}

func findUserAccountFromActiveToken(r *http.Request) (AdminAccount, bool) {
	var acc AdminAccount
	c, err := r.Cookie(defaultCookieName)
	if err == nil {
		email, found := sessionMap[c.Value]
		if found != false {
			acc = activeAccounts[email]
			return acc, true
		}
	}

	return acc, false
}

func cleanupExpiredTokens(creds Credentials) {
	account, found := activeAccounts[creds.Email]

	if found {
		shouldClean := account.hasTokenExpired()
		if shouldClean && account.SessionToken != "" {
			account.CleanToken()
		}
	}
}
