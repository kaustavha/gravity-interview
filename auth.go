package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

var sessionTokens = []string{}

type AdminAccount struct {
	gorm.Model
	Email           string
	SessionToken    string
	SessionExpiry   time.Time
	Password        string
	IsUpgraded      bool
	AccountId       string
	Users           int
	MaxUsers        int
	AssociatedUsers []Metric
}

func (a *AdminAccount) SaveInDB() {
	metricDB := GetDB()
	dbconn := metricDB.getConn()
	if dbconn.NewRecord(a) {
		dbconn.Create(&a)
	} else {
		dbconn.Update(&a)
	}
}

func (a *AdminAccount) CountAssociatedUsers() int {
	metricDB := GetDB()
	c := metricDB.countAllUniqueUsersInAccount(a.AccountId)
	return c
}

func (a *AdminAccount) hasTokenExpired() bool {
	now := time.Now()
	hasPassed := a.SessionExpiry.Before(now)
	return hasPassed
}

// func (a *AdminAccount) findSelf() *AdminAccount {
// 	metricDB := GetDB()
// 	dbconn := metricDB.getConn()
// 	dbconn.First(&a) // only 1 hardcoded admin user for now

// }

func (a *AdminAccount) CleanToken() {

	clearSessionDetails(*a)
	token := a.SessionToken
	sessionTokenIndex := index(sessionTokens, token)
	if sessionTokenIndex == -1 {
		return
	}

	sessionTokens = deleteInSlice(sessionTokens, sessionTokenIndex)

	sessionMap[a.SessionToken] = a.Email
	a.SessionToken = ""
	a.SessionExpiry = time.Now()

	fmt.Println(sessionTokens, sessionMap, activeAccounts)
	fmt.Println("cleant oken", a.SessionToken)
	// a.UpdateSelf()

	// delete(sessionMap, token)
}

func (a *AdminAccount) UpdateSelf() {
	updateSessionDetails(*a)
	// updateSessionDetails(*a)
	// fmt.Println("UpdateSelf", activeAccounts[a.Email].SessionToken)
	// activeAccounts[a.Email] = *a
	// fmt.Println("UpdateSelf", activeAccounts[a.Email].SessionToken)
	// if a.SessionToken != "" {
	// 	fmt.Println("hmm", a.SessionToken)
	// 	sessionMap[a.SessionToken] = a.Email
	// 	sessionTokens = append(sessionTokens, a.SessionToken)
	// }
}

func clearSessionDetails(a AdminAccount) {
	fmt.Println("s u", sessionTokens, sessionMap, activeAccounts)
	token := a.SessionToken
	sessionTokenIndex := index(sessionTokens, token)
	if sessionTokenIndex == -1 {
		return
	}

	sessionTokens = deleteInSlice(sessionTokens, sessionTokenIndex)

	delete(sessionMap, token)

	activeAccounts[a.Email] = a
	fmt.Println(sessionTokens, sessionMap, activeAccounts)
}

func updateSessionDetails(a AdminAccount) {
	activeAccounts[a.Email] = a
	if a.SessionToken != "" {
		sessionMap[a.SessionToken] = a.Email
		sessionTokens = append(sessionTokens, a.SessionToken)
	}
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var expected Credentials
var activeAccounts map[string]AdminAccount

var mySigningKey []byte

// map session -> email for lookup in activeAccounts map
var sessionMap map[string]string

func initAuth() {
	// sessionTokens = []string{}
	activeAccounts = make(map[string]AdminAccount)
	sessionMap = make(map[string]string)
	mySigningKey = []byte(defaultHashedPass)

	expected = Credentials{
		Email:    defaultEmail,
		Password: defaultHashedPass,
	}

}

func decodeAndCheckCreds(r *http.Request) (Credentials, int) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		fmt.Println("body decode err")
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
		fmt.Println("found cookie")
		fmt.Println(sessionTokens, sessionToken)
		return true
	}
	return false
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

	account, found := activeAccounts[creds.Email]

	expiry := time.Now().Add(120 * time.Second)
	tokenString := getJWT(creds, expiry, mySigningKey)
	if found {
		account.SessionToken = tokenString
		account.SessionExpiry = expiry
	} else {
		account = AdminAccount{
			Email:         creds.Email,
			Password:      creds.Password,
			IsUpgraded:    false,
			SessionExpiry: expiry,
			SessionToken:  tokenString,
			AccountId:     defaultAccountID,
			Users:         0,
			MaxUsers:      maxUsers,
		}
	}

	account.UpdateSelf()
	// activeAccounts[account.Email] = account
	// sessionMap[tokenString] = creds.Email
	// sessionTokens = append(sessionTokens, tokenString)

	http.SetCookie(w, &http.Cookie{
		Name:    defaultCookieName,
		Value:   account.SessionToken,
		Expires: account.SessionExpiry,
	})

	w.WriteHeader(http.StatusOK)
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
