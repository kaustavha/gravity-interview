package authenticator

import (
	"encoding/json"
	"time"

	"github.com/gravitational/trace"
	"github.com/jinzhu/gorm"
)

// AdminAccount is a main hardcoded admin
type AdminAccount struct {
	gorm.Model
	Email                string
	SessionToken         string
	SessionExpiry        time.Time
	AccountID            string
	IsUpgraded           bool
	Users                int
	MaxUsers             int
	MaxUsersAfterUpgrade int
}

// DashboardInfo is used to convert internal state for front end gets
type DashboardInfo struct {
	Users      int  `json:"Users"`
	IsUpgraded bool `json:"IsUpgraded"`
	MaxUsers   int  `json:"MaxUsers"`
}

// GetDefaultAdminAccount returns a default admin account
func GetDefaultAdminAccount(expected defaults) *AdminAccount {
	acc := &AdminAccount{
		Email:                expected.Email,
		IsUpgraded:           false,
		AccountID:            expected.AccountID,
		Users:                0,
		MaxUsers:             expected.MaxUsers,
		MaxUsersAfterUpgrade: expected.MaxUsersAfterUpgrade,
	}
	return acc
}

// HasTokenExpired checks if our auth token has expired in our cookie
func (a *AdminAccount) HasTokenExpired() bool {
	now := time.Now()
	hasPassed := a.SessionExpiry.Before(now)
	return hasPassed
}

// ClearToken empties the token on admin
func (a *AdminAccount) ClearToken() {
	a.SessionExpiry = time.Now()
	a.SessionToken = ""
}

// UpdateToken refreshes a admin accs auth token and session expire
func (a *AdminAccount) UpdateToken(password string, mySigningKey []byte) {
	expiry := time.Now().Add(120 * time.Minute)
	tokenString := getJWT(a.Email, password, expiry, mySigningKey)
	a.SessionExpiry = expiry
	a.SessionToken = tokenString
}

// Upgrade upgrades an admin acc
func (a *AdminAccount) Upgrade() {
	a.IsUpgraded = true
	a.MaxUsers = a.MaxUsersAfterUpgrade
}

// GetInfo returns json encoded acc info for frontend
func (a *AdminAccount) GetInfo() ([]byte, error) {
	res := DashboardInfo{
		Users:      a.Users,
		IsUpgraded: a.IsUpgraded,
		MaxUsers:   a.MaxUsers,
	}
	resJSON, err := json.Marshal(res)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return resJSON, nil
}

// IncrementUserCount used by the iotdata handler to increment the users in our acc
func (a *AdminAccount) IncrementUserCount() {
	a.Users++
}
