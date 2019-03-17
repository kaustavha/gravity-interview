package authenticator

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

func deleteAllInSlice(a []string, t string) []string {
	i := index(a, t)
	if i == -1 {
		return a
	}
	a = deleteInSlice(a, i)
	return deleteAllInSlice(a, t)
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

func hasTokenExpired(expiry time.Time) bool {
	now := time.Now()
	hasPassed := expiry.Before(now)
	return hasPassed
}

func getJWT(email string, password string, expiry time.Time, mySigningKey []byte) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["Email"] = email
	claims["Password"] = password
	claims["exp"] = expiry
	tokenString, _ := token.SignedString(mySigningKey)
	return tokenString
}
