package authenticator

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gravitational/trace"
	"github.com/jinzhu/gorm"
)

type defaults struct {
	AccountID            string
	Email                string
	Password             string
	MaxUsers             int
	MaxUsersAfterUpgrade int
	SigningKey           []byte
	DefaultCookieName    string
}

// Authenticator is main struct for authentication actions
type Authenticator struct {
	Database *DB
	Expected defaults
	Sessions map[string]*AdminAccount // map session -> email for lookup in activeAccounts map
	Tokens   []string                 // set of active session token strings
}

// NewAuthenticator returns a new preconfigured authenticator
func NewAuthenticator(a string, e string, p string, m int, s []byte, ma int, c string, db *gorm.DB) (*Authenticator, error) {
	authenticator := &Authenticator{
		Expected: defaults{
			AccountID:            a,
			Email:                e,
			Password:             p,
			MaxUsers:             m,
			SigningKey:           s,
			MaxUsersAfterUpgrade: ma,
			DefaultCookieName:    c,
		},
		Database: &DB{
			dbconn:    db,
			tableName: "admin_accounts",
		},
		Tokens:   []string{},
		Sessions: make(map[string]*AdminAccount),
	}
	err := authenticator.Database.Setup()
	return authenticator, err
}

//LogoutHandler clears a users session and logs them out
func (a *Authenticator) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionToken, err := a.getSessionToken(r)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	acc, err := a.logout(sessionToken)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    a.Expected.DefaultCookieName,
		Value:   acc.SessionToken,
		Expires: acc.SessionExpiry,
	})
	w.WriteHeader(http.StatusOK)
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//LoginHandler logs an admin user in and sets them in the session
func (a *Authenticator) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if a.IsAuthenticated(r) {
		fmt.Println("Already authed")
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

	acc, err := a.login(creds.Email, creds.Password)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    a.Expected.DefaultCookieName,
		Value:   acc.SessionToken,
		Expires: acc.SessionExpiry,
	})

	w.WriteHeader(http.StatusOK)
}

// Gets the admin user, and updates tokens
func (a *Authenticator) login(email string, password string) (*AdminAccount, error) {
	pass, err := a.decodeAndCheckCreds(password, email)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	acc := a.getAcc()
	acc.UpdateToken(pass, a.Expected.SigningKey)
	a.updateSessionDetails(acc)
	err = a.Database.SaveInDB(acc)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return acc, nil
}

func (a *Authenticator) logout(sessionToken string) (*AdminAccount, error) {
	acc, err := a.findUserAccountFromActiveToken(sessionToken)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	err = a.clearSessionDetails(acc)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	err = a.Database.SaveInDB(acc)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return acc, nil
}

// IsAuthenticated will return true if there is an active session with this user token
func (a *Authenticator) IsAuthenticated(r *http.Request) bool {
	sessionToken, err := a.getSessionToken(r)
	if err != nil {
		return false
	}
	return a.isAuthenticated(sessionToken)
}

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
		return "", trace.AccessDenied("Password doesnt match known admin creds")
	}

	if a.Expected.Email != email {
		return "", trace.AccessDenied("Email doesnt match known admin creds")
	}

	codedPass, err := HashPassword(password)
	if err != nil {
		return "", trace.AccessDenied("Coudnt hash pass")
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

//CleanupExpiredTokens used in middlewares to clear expired tokens on calls
func (a *Authenticator) CleanupExpiredTokens() error {
	for _, account := range a.Sessions {
		shouldClean := account.HasTokenExpired()
		if shouldClean && account.SessionToken != "" {
			err := a.clearSessionDetails(account)
			if err != nil {
				return trace.Wrap(err)
			}
		}
	}
	return nil
}

func (a *Authenticator) findUserAccountFromActiveToken(token string) (*AdminAccount, error) {
	acc, found := a.Sessions[token]
	if found != false {
		return acc, nil
	}
	return acc, trace.NotFound("RecordNotFound")
}

//Upgrade an admin user and increase user storage cap
func (a *Authenticator) Upgrade(r *http.Request) ([]byte, error) {
	token, err := a.getSessionToken(r)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	acc, err := a.findUserAccountFromActiveToken(token)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	acc.Upgrade()
	a.Database.SaveInDB(acc)
	a.updateSessionDetails(acc)
	info, err := acc.GetInfo()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return info, nil
}

//GetInfo returns the json info expected by the front end
func (a *Authenticator) GetInfo(r *http.Request) ([]byte, error) {
	token, err := a.getSessionToken(r)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	acc, err := a.findUserAccountFromActiveToken(token)
	if err != nil {
		if trace.IsNotFound(err) {
			acc, err = a.Database.FindAdmin(a.Expected.AccountID)
			a.updateSessionDetails(acc)
		}
		if err != nil {
			return nil, trace.Wrap(err)
		}
	}

	resJSON, err := acc.GetInfo()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return resJSON, nil
}

// IncrementUserCount is used by the iotdata handler to increment the users in our acc
func (a *Authenticator) IncrementUserCount(accountID string) error {
	acc, err := a.Database.FindAdmin(accountID)
	if err != nil {
		return trace.Wrap(err)
	}

	acc.IncrementUserCount()
	// a.Database.SaveInDB(acc)
	a.updateSessionDetails(acc)
	return nil
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
		return "", trace.NotFound("Cookie not found")
	}
	sessionToken := c.Value
	return sessionToken, nil
}
