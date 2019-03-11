package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	maxUsers         = 100
	maxUsersUpgraded = 1000
	defaultPort      = ":8443"

	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"

	defaultCookieName = "session_token"
	defaultSigningKey = "mysecretsigningkey"

	pemPath = "./fixtures/server-cert.pem"
	keyPath = "./fixtures/server-key.pem"

	defaultHashedPass = "$2a$14$JMgUM09OV3HPAMKNM9nnb.wghzq5ayYRe91li1j9uqc9pGxU0kQX2"
	defaultEmail      = "a@a.com"
)

// Here we are implementing the NotImplemented handler. Whenever an API endpoint is hit
// we will simply return the message "Not Implemented"
var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

func UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	acc, found := findUserAccountFromActiveToken(r)
	if !found {
		w.WriteHeader(http.StatusNotFound)
	}
	if acc.IsUpgraded == true {
		w.WriteHeader(http.StatusLoopDetected)
		return
	}
	acc.IsUpgraded = true
	acc.MaxUsers = maxUsersUpgraded
	// updateSessionDetails(acc)
	acc.UpdateSelf()
	acc.SaveInDB()
	onSuccesfulUpgrade(w, acc)
}

func UpgradeCheckHandler(w http.ResponseWriter, r *http.Request) {
	acc, found := findUserAccountFromActiveToken(r)

	// Get latest state of acc
	if !found {
		w.WriteHeader(http.StatusNotFound)
	}
	// force update from db
	foundInDB, dbacc := db.findAdmin(acc.AccountId)
	if foundInDB {
		acc = *dbacc
	}
	// acc.UpdateSelf()
	onSuccesfulUpgrade(w, acc)
}

func onSuccesfulUpgrade(w http.ResponseWriter, acc AdminAccount) {
	resJSON, err := json.Marshal(acc)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(resJSON)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if isAuthenticated(r) {
			next(w, r)
			return
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RequestURI)
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

func main() {
	fmt.Println("App boot up...")
	// Create global vars
	InitAuth()
	createDBConn()

	WrappedLoginHandler := loggingMiddleware(cleanupExpiredTokensMiddleware(LoginHandler))
	WrappedIOTDataHandler := loggingMiddleware(IOTDataHandler)

	http.HandleFunc("/api/login", WrappedLoginHandler)
	http.HandleFunc("/api/logout", applyMiddlewares(LogoutHandler))
	http.HandleFunc("/api/authcheck", applyMiddlewares(AuthcheckHandler))

	http.HandleFunc("/api/dashboard", applyMiddlewares(DashboardHandler))

	http.HandleFunc("/api/upgrade", applyMiddlewares(UpgradeHandler))
	http.HandleFunc("/api/upgradecheck", applyMiddlewares(UpgradeCheckHandler))

	http.HandleFunc("/metrics", WrappedIOTDataHandler)

	port := os.Getenv("PORT") // TODO - add env vars
	if port == "" {
		port = defaultPort
	}

	fmt.Println(port)

	err := http.ListenAndServeTLS(port, pemPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
