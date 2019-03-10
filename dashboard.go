package main

import (
	"encoding/json"
	"net/http"
)

type DashboardInfo struct {
	UserCount int `json:"userCount"`
}

func countUsers(accId string) int {
	metricDB := GetDB()
	c := metricDB.countAllUniqueUsersInAccount(accId)
	return c
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	acc := findUserAccountFromActiveToken(r)
	acc.Users = countUsers(acc.AccountId)
	setAccountInfo(acc)

	dasboardInfo := DashboardInfo{
		UserCount: acc.Users,
	}

	w.Header().Set(contentTypeHeader, contentTypeJSON)
	w.WriteHeader(http.StatusOK)
	resJSON, err := json.Marshal(dasboardInfo)
	if err == nil {
		w.Write(resJSON)
	}
}
