package main

import (
	"encoding/json"
	"net/http"
)

func UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	acc, found := findUserAccountFromActiveToken(r)
	if !found {
		w.WriteHeader(http.StatusNotFound)
	}
	if acc.IsUpgraded == true {
		w.WriteHeader(http.StatusLoopDetected)
		return
	}
	acc.IsUpgraded = true
	acc.MaxUsers = maxUsersUpgraded
	acc.UpdateSelf()
	acc.SaveInDB()
	onSuccesfulUpgrade(w, acc)
}

func UpgradeCheckHandler(w http.ResponseWriter, r *http.Request) {
	acc, found := findUserAccountFromActiveToken(r)

	// Get latest state of acc
	if !found {
		w.WriteHeader(http.StatusNotFound)
	}
	// force update from db
	foundInDB, dbacc := db.findAdmin(acc.AccountId)
	if foundInDB {
		acc = *dbacc
	}
	onSuccesfulUpgrade(w, acc)
}

func onSuccesfulUpgrade(w http.ResponseWriter, acc AdminAccount) {
	resJSON, err := json.Marshal(acc)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write(resJSON)
	}
}
