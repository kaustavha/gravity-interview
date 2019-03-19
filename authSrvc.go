package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gravitational/trace"
)

//AuthService struct
type AuthService struct {
	a                 Authenticator
	defaultCookieName string
}

//GetNewAuthService  returns a new instance of DashboardService
func GetNewAuthService(a interface{}, cookieName string) *AuthService {
	return &AuthService{
		a:                 reflect.ValueOf(a).Interface().(Authenticator),
		defaultCookieName: cookieName,
	}
}

//AuthcheckHandler basic handler, returns ok always, auth check is done by middleware
func (as *AuthService) AuthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

//LogoutHandler clears a users session and logs them out
func (as *AuthService) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token, expiry, err := as.a.LogoutAdmin(r)
	if err != nil {
		fmt.Println(err.Error())
		if trace.IsNotFound(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    as.defaultCookieName,
		Value:   token,
		Expires: expiry,
	})
	w.WriteHeader(http.StatusOK)
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//LoginHandler logs an admin user in and sets them in the session
func (as *AuthService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if as.a.IsAuthenticated(r) {
		w.WriteHeader(http.StatusOK)
		return
	}

	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, expiry, err := as.a.LoginAdmin(creds.Email, creds.Password)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    as.defaultCookieName,
		Value:   token,
		Expires: expiry,
	})

	w.WriteHeader(http.StatusOK)
}
