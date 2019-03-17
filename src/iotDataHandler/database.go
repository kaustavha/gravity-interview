package iotdatahandler

import (
	"github.com/gravitational/trace"
	"github.com/jinzhu/gorm"
)

//IotDataHandlerDB is the main struct
type IotDataHandlerDB struct {
	dbconn    *gorm.DB
	tableName string
}

//GetNewIotDataHandlerDB returns a new IotDataHandlerDB
func GetNewIotDataHandlerDB(db *gorm.DB) *IotDataHandlerDB {
	db.AutoMigrate(&Metric{})
	return &IotDataHandlerDB{
		dbconn:    db,
		tableName: "metrics",
	}
}

//SaveInDB saves a metric in the users db and updates the user count in the assoc admin
func (db *IotDataHandlerDB) SaveInDB(m *Metric) error {
	_, err := db.findMetric(m)

	// handle duplciate users
	if trace.IsNotFound(err) {
		err = db.dbconn.Create(&m).Error
	} else if err == nil {
		err = db.updateMetric(m)
	}

	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (db *IotDataHandlerDB) updateMetric(m *Metric) error {
	err := db.dbconn.Table(db.tableName).Where(&Metric{
		AccountID: m.AccountID,
		UserID:    m.UserID,
	}).Updates(m)

	if err.Error != nil {
		return trace.Wrap(err.Error)
	}
	return nil
}

func (db *IotDataHandlerDB) findMetric(m *Metric) (*Metric, error) {
	metric := &Metric{}
	record := db.dbconn.Table(db.tableName).Where(&Metric{
		AccountID: m.AccountID,
		UserID:    m.UserID,
	}).Find(&metric)
	if record.RecordNotFound() {
		return nil, trace.NotFound("RecordNotFound")
	}

	// All unahndled errors e.g. db conn errs
	if record.Error != nil {
		return nil, trace.Wrap(record.Error)
	}
	return metric, nil
}
