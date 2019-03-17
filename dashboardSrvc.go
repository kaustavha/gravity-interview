package main

import (
	"fmt"
	"net/http"
	"reflect"
)

//DashboardService struct
type DashboardService struct {
	a Authenticator
}

//GetNewDashboardService  returns a new instance of DashboardService
func GetNewDashboardService(a interface{}) *DashboardService {
	return &DashboardService{
		a: reflect.ValueOf(a).Interface().(Authenticator),
	}
}

//DashboardHandler Http handler for dashboard info requests
func (d *DashboardService) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	d.getInfo(w, r)
}

//UpgradeCheckHandler Http handler for dashboard info requests
func (d *DashboardService) UpgradeCheckHandler(w http.ResponseWriter, r *http.Request) {
	d.getInfo(w, r)
}

//UpgradeHandler http handler for upgrading admin users
func (d *DashboardService) UpgradeHandler(w http.ResponseWriter, r *http.Request) {
	resJSON, err := d.a.Upgrade(r)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}

func (d *DashboardService) getInfo(w http.ResponseWriter, r *http.Request) {
	resJSON, err := d.a.GetInfo(r)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resJSON)
}
