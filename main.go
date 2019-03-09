package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func APIHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("api hit")
	w.WriteHeader(http.StatusOK)
}

// Define our struct
type authenticationMiddleware struct {
	tokenUsers map[string]string
}

// Initialize it somewhere
func (amw *authenticationMiddleware) Populate() {
	amw.tokenUsers = make(map[string]string)
	amw.tokenUsers["00000000"] = "user0"
}

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		token := r.Header.Get("X-Session-Token")

		fmt.Println(token, r.Header, amw.tokenUsers)
		if user, found := amw.tokenUsers[token]; found {
			// We found the token in our map
			log.Printf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func main() {
	fmt.Println("App boot up...")
	router := mux.NewRouter()
	router.HandleFunc("/api", APIHandler)

	// Set up auth middleware
	amw := authenticationMiddleware{}
	amw.Populate()
	router.Use(amw.Middleware)

	port := os.Getenv("PORT") // TODO
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
