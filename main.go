package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

const (
	maxUsers         = 100
	maxUsersUpgraded = 1000
)

// Here we are implementing the NotImplemented handler. Whenever an API endpoint is hit
// we will simply return the message "Not Implemented"
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	acc := findUserAccountFromActiveToken(r)
	cleanToken(acc)
	w.WriteHeader(http.StatusOK)
	onSuccesfulUpgrade(w, acc)
}

func UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	acc := findUserAccountFromActiveToken(r)
	if acc.IsUpgraded == true {
		w.WriteHeader(http.StatusLoopDetected)
		return
	}
	acc.IsUpgraded = true
	acc.MaxUsers = maxUsersUpgraded
	onSuccesfulUpgrade(w, acc)
}

func UpgradeCheckHandler(w http.ResponseWriter, r *http.Request) {
	acc := findUserAccountFromActiveToken(r)
	if acc.IsUpgraded {
		onSuccesfulUpgrade(w, acc)
		return
	}
	w.WriteHeader(http.StatusUnprocessableEntity)
}

func onSuccesfulUpgrade(w http.ResponseWriter, acc AdminAccount) {
	setAccountInfo(acc)
	resJSON, err := json.Marshal(acc)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(resJSON)
	}
}

func IOTDataHandler(w http.ResponseWriter, r *http.Request) {
	NotImplemented(w, r)
}

type DashboardInfo struct {
	UserCount int `json:"userCount"`
}

// var dasboardInfo = DashboardInfo{}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	acc := findUserAccountFromActiveToken(r)
	acc.Users += 10
	setAccountInfo(acc)
	fmt.Println("DashboardHandler")

	dasboardInfo := DashboardInfo{
		UserCount: acc.Users,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resJSON, err := json.Marshal(dasboardInfo)
	if err == nil {
		w.Write(resJSON)
	}
}

var sessionTokens []string

var mySigningKey = []byte("secret")

type AdminAccount struct {
	Email         string
	SessionToken  string
	SessionExpiry time.Time
	Password      string
	IsUpgraded    bool
	AccountId     string
	Users         int
	MaxUsers      int
}

var activeAccounts map[string]AdminAccount

// map session -> email for lookup in activeAccounts map
var sessionMap map[string]string

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var expected = Credentials{
	Email:    "a@a.com",
	Password: "$2a$14$JMgUM09OV3HPAMKNM9nnb.wghzq5ayYRe91li1j9uqc9pGxU0kQX2",
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

// StatusOK if already signed in, or after signin based on incoming cookie or email+password
// otherwise frontend needs to prompt for usrname and pass and try again
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	creds, status := decodeAndCheckCreds(r)
	if status != http.StatusOK {
		w.WriteHeader(status)
	}

	account, found := activeAccounts[creds.Email]

	expiry := time.Now().Add(120 * time.Second)
	tokenString := getJWT(creds, expiry)
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
			AccountId:     uuid.New(),
			Users:         0,
			MaxUsers:      maxUsers,
		}
	}

	activeAccounts[account.Email] = account
	sessionMap[tokenString] = creds.Email
	sessionTokens = append(sessionTokens, tokenString)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   account.SessionToken,
		Expires: account.SessionExpiry,
	})

	w.WriteHeader(http.StatusOK)
}

func getJWT(creds Credentials, expiry time.Time) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["Email"] = creds.Email
	claims["Password"] = creds.Password
	claims["exp"] = expiry
	tokenString, _ := token.SignedString(mySigningKey)
	return tokenString
}

func setAccountInfo(acc AdminAccount) {
	activeAccounts[acc.Email] = acc
}

func findUserAccountFromActiveToken(r *http.Request) AdminAccount {
	var acc AdminAccount
	c, err := r.Cookie("session_token")
	if err == nil {
		email, found := sessionMap[c.Value]
		if found != false {
			acc = activeAccounts[email]
			return acc
		}
	}

	return acc
}

func isAuthenticated(r *http.Request) bool {
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Println("No Cookie")
			return false
		}
		fmt.Println("err", err)
		return false
	}
	sessionToken := c.Value
	found := index(sessionTokens, sessionToken)
	if found != -1 {
		return true
	}
	fmt.Println("not found token", sessionTokens, sessionToken)
	return false
}

func index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}

func deleteInSlice(a []string, i int) []string {
	a = append(a[:i], a[i+1:]...)
	return a
}

func cleanupExpiredTokens(creds Credentials) {
	account, found := activeAccounts[creds.Email]

	if found {
		shouldClean := hasTokenExpired(account.SessionExpiry)
		if shouldClean {
			cleanToken(account)
		}
	}
}
func cleanToken(acc AdminAccount) {
	token := acc.SessionToken
	sessionTokenIndex := index(sessionTokens, token)
	if sessionTokenIndex != -1 {
		sessionTokens = deleteInSlice(sessionTokens, sessionTokenIndex)
	}

	acc.SessionToken = ""
	acc.SessionExpiry = time.Now()
	setAccountInfo(acc)

	delete(sessionMap, token)
}

func hasTokenExpired(expiry time.Time) bool {
	now := time.Now()
	hasPassed := expiry.Before(now)
	return hasPassed
}

func APIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("api hit")
	w.WriteHeader(http.StatusOK)
}

func AuthcheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("AuthcheckHandler hit")
	w.WriteHeader(http.StatusOK)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only login route doesnt need auth
		dest := r.RequestURI
		fmt.Println(dest)
		if isAuthenticated(r) {
			if dest == "/api/login" {
				w.WriteHeader(http.StatusOK)
			} else {
				next.ServeHTTP(w, r)
			}
			return
		}

		if dest == "/api/login" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func cleanupExpiredTokensMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanupExpiredTokens(expected)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Create global vars
	activeAccounts = make(map[string]AdminAccount)
	sessionMap = make(map[string]string)

	fmt.Println("App boot up...")
	router := mux.NewRouter()

	router.Use(authMiddleware)
	router.Use(cleanupExpiredTokensMiddleware)

	router.HandleFunc("/api", APIHandler)
	router.HandleFunc("/api/login", LoginHandler)
	router.HandleFunc("/api/logout", LogoutHandler)
	router.HandleFunc("/api/authcheck", AuthcheckHandler)

	router.HandleFunc("/api/dashboard", DashboardHandler)

	router.HandleFunc("/api/upgrade", UpgradeHandler)
	router.HandleFunc("/api/upgradecheck", UpgradeCheckHandler)

	router.HandleFunc("/metrics", IOTDataHandler)

	port := os.Getenv("PORT") // TODO - add env vars
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
