package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

const (
	defaultBearerToken = "shmoken"
	defaultAccountID   = "5a28fa21-c70d-4bf3-b4c4-c4b109d5d269"
)

// Metric is a metric send by the iot device
// every time user logs into it
type Metric struct {
	gorm.Model
	// AccountID is a unique UUID identifying the account
	AccountID string `json:"account_id"`
	// UserID is a unique ID identityfing the user
	// activity
	UserID string `json:"user_id"`
	// Timestamp is a time as recorded by the device
	Timestamp time.Time `json:"timestamp"`
}

func (m *Metric) SaveInDB() {
	metricDB := GetDB()
	dbconn := metricDB.getConn()
	dbconn.Create(&m)
}

// String returns debug-friendly representation of the metric
func (m *Metric) String() string {
	data, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// IOTDataHandler Handle data coming from iot data gen
func IOTDataHandler(w http.ResponseWriter, r *http.Request) {
	// Example incoming
	// {
	// 	"account_id": "781df840-09da-42f4-ba29-996d2ff76a73",
	// 	"user_id": "bf506b23-8c4e-4c8e-af95-e331dba766ab",
	// 	"timestamp": "2019-03-03T18:02:30.424878129Z"
	//   }
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) == 1 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("token fail"))
		return
	}

	reqToken = splitToken[1]
	if reqToken != defaultBearerToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	contentType := r.Header.Get(contentTypeHeader)
	if contentTypeJSON != contentType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var metric Metric

	err := json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(metric.AccountID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if len(metric.UserID) == 0 {
		fmt.Println("user id fail", err, r, r.Body, r.Header)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metric.SaveInDB()

	w.WriteHeader(http.StatusOK)
}
