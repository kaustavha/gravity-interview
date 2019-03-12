package authenticator

import (
	"time"

	"github.com/gravitational/trace"
	"github.com/jinzhu/gorm"
)

type Defaults struct {
	AccountID  string
	Email      string
	Password   string
	MaxUsers   int
	SigningKey []byte
	TableName  string //admin_accounts
}

type Authenticator struct {
	Database       *DB
	Expected       Defaults
	Sessions       map[string]string        // map session -> email for lookup in activeAccounts map
	Tokens         []string                 // set of active session token strings
	ActiveAccounts map[string]*AdminAccount // map of email -> Adminaccs
}

func NewAuthenticator(a string, e string, p string, m int, s []byte, db *gorm.DB) (*Authenticator, error) {
	authenticator := &Authenticator{
		Expected: Defaults{
			AccountID:  a,
			Email:      e,
			Password:   p,
			MaxUsers:   m,
			SigningKey: s,
			TableName:  "admin_accounts",
		},
		Database: &DB{
			dbconn: db,
		},
		Tokens:         []string{},
		ActiveAccounts: make(map[string]*AdminAccount),
		Sessions:       make(map[string]string),
	}

	err := authenticator.Database.Setup(authenticator.Expected.TableName)

	return authenticator, err
}

// Gets the admin user, and updates tokens
func (a *Authenticator) Login(sessionToken string, password string) *AdminAccount {

	acc, found := a.FindUserAccountFromActiveToken(sessionToken)

	var err error
	if found == false {
		acc, err = a.Database.FindAdmin(a.Expected.AccountID)
		if err != nil {
			acc = GetDefaultAdminAccount(a.Expected)
		}
	}
	acc.UpdateToken(password, a.Expected.SigningKey)
	a.UpdateSessionDetails(acc)
	return acc
}

// Given a un-encoded email and pass from the frontend - checks if its our defaults and returns the coded pass
func (a *Authenticator) DecodeAndCheckCreds(password string, email string) (string, error) {
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

func (a *Authenticator) clearSessionDetails(admin *AdminAccount) {
	token := admin.SessionToken
	email := admin.Email

	a.Tokens = deleteAllInSlice(a.Tokens, token)
	delete(a.Sessions, token)

	admin.SessionToken = ""
	admin.SessionExpiry = time.Now()
	a.ActiveAccounts[email] = admin
}

func (a *Authenticator) UpdateSessionDetails(admin *AdminAccount) {
	token := admin.SessionToken
	email := admin.Email
	a.ActiveAccounts[email] = admin
	if token != "" {
		a.Sessions[token] = email
		a.Tokens = append(a.Tokens, token)
	}
}

func (a *Authenticator) IsAuthenticated(sessionToken *string) bool {
	found := index(a.Tokens, *sessionToken)
	if found != -1 {
		return true
	}
	return false
}

func (a *Authenticator) CleanupExpiredTokens() {
	account, found := a.ActiveAccounts[a.Expected.Email]
	if found {
		shouldClean := account.HasTokenExpired()
		if shouldClean && account.SessionToken != "" {
			a.clearSessionDetails(account)
		}
	}
}

func (a *Authenticator) FindUserAccountFromActiveToken(token string) (*AdminAccount, bool) {
	email, found := a.Sessions[token]
	var acc *AdminAccount
	if found != false {
		acc = a.ActiveAccounts[email]
		return acc, true
	}
	return acc, false
}
