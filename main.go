package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gravitational/trace"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kaustavha/gravity-interview/src/authenticator"
	"github.com/kaustavha/gravity-interview/src/iotDataHandler"
)

const (
	defaultPort = ":8443"

	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"

	defaultCookieName = "session_token"

	pemPath = "./fixtures/server-cert.pem"
	keyPath = "./fixtures/server-key.pem"

	accountID        = "5a28fa21-c70d-4bf3-b4c4-c4b109d5d269"
	email            = "a@a.com"
	hashedPass       = "$2a$14$JMgUM09OV3HPAMKNM9nnb.wghzq5ayYRe91li1j9uqc9pGxU0kQX2"
	maxUsers         = 100
	maxUsersUpgraded = 1000
	signingKey       = "bXlzZWNyZXRzaWduaW5na2V5Cg=="

	debugDefault = true
)

func createDBConn() (*gorm.DB, error) {
	const (
		dbhost    = "localhost"
		dbport    = "5432"
		dbuser    = "postgres"
		dbname    = "iotdb"
		dbpass    = "bXlzcWxwYXNzd29yZAo="
		dbsslmode = "disable"
	)

	optString := "host=" + dbhost + " " +
		"port=" + dbport + " " +
		"user=" + dbuser + " " +
		"dbname=" + dbname + " " +
		"password=" + dbpass + " " +
		"sslmode=" + dbsslmode

	conn, err := gorm.Open("postgres", optString)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return conn, nil
}

//AuthcheckHandler basic handler, returns ok always, auth check is done by middleware
func AuthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	fmt.Println("App boot up...")

	// set up err handling
	debug := os.Getenv("DEBUG")
	if debug == "" {
		trace.SetDebug(debugDefault)
	} else if debug == "FALSE" {
		trace.SetDebug(false)
	}

	db, err := createDBConn()
	if err != nil {
		trace.Fatalf("Error connecting to Database %v", err)
		fmt.Println(err.Error())
	}
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
	if err != nil {
		trace.Fatalf("Error setting up Database in Authenticator %v", err)
		fmt.Println(err.Error())
	}

	// init middleware manager
	m := GetNewMiddlewareManager(a)

	// init iotdata handler
	i := iotDataHandler.GetNewIOTDataHandler(
		a,
		contentTypeHeader,
		contentTypeJSON,
		db,
	)

	d := GetNewDashboardService(a)

	http.HandleFunc("/api/login", m.getWrappedLoginHandler(a.LoginHandler))
	http.HandleFunc("/api/authcheck", m.applyMiddlewares(AuthcheckHandler))
	http.HandleFunc("/api/logout", m.applyMiddlewares(a.LogoutHandler))
	http.HandleFunc("/metrics", m.getWrappedIOTDataHandler(i.IOTDataHandler))

	http.HandleFunc("/api/dashboard", m.applyMiddlewares(d.DashboardHandler))
	http.HandleFunc("/api/upgrade", m.applyMiddlewares(d.UpgradeHandler))
	http.HandleFunc("/api/upgradecheck", m.applyMiddlewares(d.UpgradeCheckHandler))

	port := os.Getenv("PORT") // TODO - add env vars
	if port == "" {
		port = defaultPort
	}

	fmt.Println("Listening on: ", port)

	err = http.ListenAndServeTLS(port, pemPath, keyPath, nil)
	if err != nil {
		trace.Fatalf("ListenAndServe: %v", err)
	}
}
