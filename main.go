package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gravitational/trace"
	"github.com/jinzhu/gorm"
	"github.com/kaustavha/gravity-interview/src/authenticator"
)

const (
	defaultPort = ":8443"

	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"

	defaultCookieName = "session_token"

	pemPath = "./fixtures/server-cert.pem"
	keyPath = "./fixtures/server-key.pem"

	AccountID        = "5a28fa21-c70d-4bf3-b4c4-c4b109d5d269"
	Email            = "a@a.com"
	HashedPass       = "$2a$14$JMgUM09OV3HPAMKNM9nnb.wghzq5ayYRe91li1j9uqc9pGxU0kQX2"
	maxUsers         = 100
	maxUsersUpgraded = 1000
	SigningKey       = "bXlzZWNyZXRzaWduaW5na2V5Cg=="
	TableName        = "metrics"
)

func createDBConn() (*gorm.DB, error) {

	optString := "host=" + dbhost + " " +
		"port=" + dbport + " " +
		"user=" + dbuser + " " +
		"dbname=" + dbname + " " +
		"password=" + dbpass + " " +
		"sslmode=" + dbsslmode

	fmt.Println(optString)

	conn, err := gorm.Open("postgres", optString)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return conn, nil
}

func main() {
	fmt.Println("App boot up...")
	db, err := createDBConn()
	if err != nil {
		trace.Wrap(err)
	}
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
	if err != nil {
		trace.Wrap(err)
	}

	// init middleware manager
	m := GetNewMiddlewareManager(a)

	// loggingMiddleware(ExampleLoginHandler)

	// fmt.Println(a)
	// Create global vars
	InitAuth()
	CreateDBConn(db)
	defer GetDBConn().Close()
	InitRoutes()

	http.HandleFunc("/api/login", m.loggingMiddleware(a.LoginHandler))
	http.HandleFunc("/api/authcheck", m.applyMiddlewares(AuthcheckHandler))
	http.HandleFunc("/api/logout", m.applyMiddlewares(a.LogoutHandler))

	port := os.Getenv("PORT") // TODO - add env vars
	if port == "" {
		port = defaultPort
	}

	fmt.Println("Listening on: ", port)

	err = http.ListenAndServeTLS(port, pemPath, keyPath, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
