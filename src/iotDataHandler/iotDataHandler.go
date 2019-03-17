package iotdatahandler

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gravitational/trace"
	"github.com/jinzhu/gorm"
)

type Authenticator interface {
	IncrementUserCount(AccountID string) error
}

type IOTDataHandler struct {
	db       *IotDataHandlerDB
	a        Authenticator
	defaults *defaults
}

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

type defaults struct {
	contentTypeHeader  string
	contentTypeJSON    string
	defaultAccountID   string
	defaultBearerToken string
}

//GetNewIOTDataHandler returns a new isntance of our iotdata handler
func GetNewIOTDataHandler(a interface{}, contentTypeHeader string, contentTypeJSON string, defaultAccountID string, defaultBearerToken string, db *gorm.DB) *IOTDataHandler {
	return &IOTDataHandler{
		a: reflect.ValueOf(a).Interface().(Authenticator),
		defaults: &defaults{
			contentTypeHeader:  contentTypeHeader,
			contentTypeJSON:    contentTypeJSON,
			defaultAccountID:   defaultAccountID,
			defaultBearerToken: defaultBearerToken,
		},
		db: GetNewIotDataHandlerDB(db),
	}
}

//SaveInDB saves metrics in a metrics table and also increments the assoc admin users usercount
func (i *IOTDataHandler) SaveInDB(m *Metric) error {
	err := i.db.SaveInDB(m)
	if err != nil {
		return trace.Wrap(err)
	}
	err = i.a.IncrementUserCount(m.AccountID)
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
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
func (i *IOTDataHandler) IOTDataHandler(w http.ResponseWriter, r *http.Request) {
	// Example incoming
	// {
	// 	"account_id": "781df840-09da-42f4-ba29-996d2ff76a73",
	// 	"user_id": "bf506b23-8c4e-4c8e-af95-e331dba766ab",
	// 	"timestamp": "2019-03-03T18:02:30.424878129Z"
	//   }
	contentTypeHeader := i.defaults.contentTypeHeader
	contentTypeJSON := i.defaults.contentTypeJSON
	w.Header().Set(contentTypeHeader, contentTypeJSON)
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")

	if len(splitToken) == 1 {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(trace.AccessDenied("Incoming token is empty").Error()))
		return
	}

	reqToken = splitToken[1]
	if reqToken != i.defaults.defaultBearerToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	contentType := r.Header.Get(contentTypeHeader)
	if contentTypeJSON != contentType {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var metric *Metric

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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// untested case from data generator tests
	if time.Time.IsZero(metric.Timestamp) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	i.SaveInDB(metric)
	w.WriteHeader(http.StatusOK)
}
