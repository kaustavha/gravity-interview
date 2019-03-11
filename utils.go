package main

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

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

func hasTokenExpired(expiry time.Time) bool {
	now := time.Now()
	hasPassed := expiry.Before(now)
	return hasPassed
}

func getJWT(creds Credentials, expiry time.Time, mySigningKey []byte) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["Email"] = creds.Email
	claims["Password"] = creds.Password
	claims["exp"] = expiry
	tokenString, _ := token.SignedString(mySigningKey)
	return tokenString
}
