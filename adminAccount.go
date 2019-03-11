package main

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

func (a *AdminAccount) SaveInDB() {
	metricDB := GetDB()
	dbconn := metricDB.getConn()
	found, _ := metricDB.findAdmin(a.AccountId)
	if dbconn.NewRecord(a) && found != true {
		dbconn.Create(&a)
	} else {
		metricDB.updateById(*a)
	}
}

func (a *AdminAccount) CountAssociatedUsers() int {
	metricDB := GetDB()
	c := metricDB.countAllUniqueUsersInAccount(a.AccountId)
	return c
}

func (a *AdminAccount) hasTokenExpired() bool {
	now := time.Now()
	hasPassed := a.SessionExpiry.Before(now)
	return hasPassed
}

func (a *AdminAccount) CleanToken() {
	clearSessionDetails(*a)
}

func (a *AdminAccount) UpdateSelf() {
	updateSessionDetails(*a)
}

func clearSessionDetails(a AdminAccount) {
	token := a.SessionToken
	sessionTokens = deleteAllInSlice(sessionTokens, token)
	delete(sessionMap, token)

	a.SessionToken = ""
	a.SessionExpiry = time.Now()
	activeAccounts[a.Email] = a
}

func updateSessionDetails(a AdminAccount) {
	activeAccounts[a.Email] = a
	if a.SessionToken != "" {
		sessionMap[a.SessionToken] = a.Email
		sessionTokens = append(sessionTokens, a.SessionToken)
	}
}
