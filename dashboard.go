package main

import (
	"encoding/json"
	"net/http"
)

type DashboardInfo struct {
	UserCount int `json:"userCount"`
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	acc, found := findUserAccountFromActiveToken(r)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if acc.Users != acc.CountAssociatedUsers() {
		acc.Users = acc.CountAssociatedUsers()
		acc.UpdateSelf()
		acc.SaveInDB()
	}

	dasboardInfo := DashboardInfo{
		UserCount: acc.Users,
	}

	w.Header().Set(contentTypeHeader, contentTypeJSON)
	resJSON, err := json.Marshal(dasboardInfo)
	if err == nil {
		w.Write(resJSON)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
}
