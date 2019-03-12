package authenticator

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AdminAccount struct {
	gorm.Model
	Email         string
	SessionToken  string
	SessionExpiry time.Time
	IsUpgraded    bool
	AccountId     string
	Users         int
	MaxUsers      int
}

func GetDefaultAdminAccount(expected Defaults) *AdminAccount {
	acc := &AdminAccount{
		Email:      expected.Email,
		IsUpgraded: false,
		AccountId:  expected.AccountID,
		Users:      0,
		MaxUsers:   expected.MaxUsers,
	}
	return acc
}

func (a *AdminAccount) HasTokenExpired() bool {
	now := time.Now()
	hasPassed := a.SessionExpiry.Before(now)
	return hasPassed
}

func (a *AdminAccount) UpdateToken(password string, mySigningKey []byte) {
	expiry := time.Now().Add(120 * time.Second)
	tokenString := getJWT(a.Email, password, expiry, mySigningKey)
	a.SessionExpiry = expiry
	a.SessionToken = tokenString
}
