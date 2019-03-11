package main

import (
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
	defaultSigningKey = "bXlzZWNyZXRzaWduaW5na2V5Cg=="

	pemPath = "./fixtures/server-cert.pem"
	keyPath = "./fixtures/server-key.pem"

	defaultHashedPass = "$2a$14$JMgUM09OV3HPAMKNM9nnb.wghzq5ayYRe91li1j9uqc9pGxU0kQX2"
	defaultEmail      = "a@a.com"
)

func main() {
	fmt.Println("App boot up...")
	// Create global vars
	InitAuth()
	CreateDBConn()
	defer GetDBConn().Close()
	InitRoutes()

	port := os.Getenv("PORT") // TODO - add env vars
	if port == "" {
		port = defaultPort
	}

	fmt.Println("Listening on: ", port)

	err := http.ListenAndServeTLS(port, pemPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
