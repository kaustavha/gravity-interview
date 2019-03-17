package authenticator

import (
	"net/http"
	"time"

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

//LogoutAdmin will logout the specified admin user
func (a *Authenticator) LogoutAdmin(r *http.Request) (string, time.Time, error) {
	sessionToken, err := a.getSessionToken(r)
	defaultTime := time.Now()
	if err != nil {
		return "", defaultTime, trace.NotFound("Could not find session token in req: %v", err)
	}
	acc, err := a.findUserAccountFromActiveToken(sessionToken)
	if err != nil {
		return "", defaultTime, trace.Wrap(err)
	}
	err = a.clearSessionDetails(acc)
	if err != nil {
		return "", defaultTime, trace.Wrap(err)
	}
	return acc.SessionToken, acc.SessionExpiry, nil
}

//LoginAdmin Gets the admin user, and updates tokens
func (a *Authenticator) LoginAdmin(email string, password string) (string, time.Time, error) {
	pass, err := a.decodeAndCheckCreds(password, email)
	defaultTime := time.Now()
	if err != nil {
		return "", defaultTime, trace.Wrap(err)
	}
	acc := a.getAcc()
	acc.UpdateToken(pass, a.Expected.SigningKey)
	a.updateSessionDetails(acc)

	err = a.Database.SaveInDB(acc)
	if err != nil {
		return "", defaultTime, trace.Wrap(err)
	}
	return acc.SessionToken, acc.SessionExpiry, nil
}

// IsAuthenticated will return true if there is an active session with this user token
func (a *Authenticator) IsAuthenticated(r *http.Request) bool {
	sessionToken, err := a.getSessionToken(r)
	if err != nil {
		return false
	}
	return a.isAuthenticated(sessionToken)
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

//GetInfo returns the json info expected by the front end fresh from our db
func (a *Authenticator) GetInfo(r *http.Request) ([]byte, error) {
	token, err := a.getSessionToken(r)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	acc, err := a.findUserAccountFromActiveToken(token)

	// Handle non-active users// overlap betw cleanup expired tokens and getinfo
	if err != nil {
		acc.AccountID = a.Expected.AccountID
	}

	acc, err = a.Database.FindAdmin(acc.AccountID)
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
	a.Database.UpdateUserCountByID(acc)
	return nil
}
