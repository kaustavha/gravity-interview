package authenticator

import (
	"net/http"

	"github.com/gravitational/trace"
)

func (a *Authenticator) isAuthenticated(sessionToken string) bool {
	found := index(a.Tokens, sessionToken)
	if found != -1 {
		return true
	}
	return false
}

// Given a un-encoded email and pass from the frontend - checks if its our defaults and returns the coded pass
func (a *Authenticator) decodeAndCheckCreds(password string, email string) (string, error) {
	if CheckPasswordHash(password, a.Expected.Password) == false {
		return "", trace.AccessDenied("password doesnt match known admin creds")
	}

	if a.Expected.Email != email {
		return "", trace.AccessDenied("email doesnt match known admin creds")
	}

	codedPass, err := HashPassword(password)
	if err != nil {
		return "", trace.AccessDenied("coudnt hash pass")
	}

	return codedPass, nil
}

func (a *Authenticator) clearSessionDetails(admin *AdminAccount) error {
	token := admin.SessionToken
	a.Tokens = deleteAllInSlice(a.Tokens, token)
	delete(a.Sessions, token)
	admin.ClearToken()
	err := a.Database.SaveInDB(admin)
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (a *Authenticator) updateSessionDetails(admin *AdminAccount) {
	token := admin.SessionToken
	if token != "" {
		a.Sessions[token] = admin
		a.Tokens = append(a.Tokens, token)
	}
}

func (a *Authenticator) findUserAccountFromActiveToken(token string) (*AdminAccount, error) {
	acc, found := a.Sessions[token]
	if found != false {
		return acc, nil
	}
	return acc, trace.NotFound("recordNotFound")
}

func (a *Authenticator) getAcc() *AdminAccount {
	acc, err := a.Database.FindAdmin(a.Expected.AccountID)
	if err != nil {
		acc = GetDefaultAdminAccount(a.Expected)
		a.Database.SaveInDB(acc)
	}

	return acc
}

func (a *Authenticator) getSessionToken(r *http.Request) (string, error) {
	c, err := r.Cookie(a.Expected.DefaultCookieName)
	if err != nil {
		return "", trace.NotFound("cookie not found")
	}
	sessionToken := c.Value
	return sessionToken, nil
}
