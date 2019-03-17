package authenticator

import (
	"github.com/gravitational/trace"
	"github.com/jinzhu/gorm"
)

//DB is a wrapper struct around gorm db
type DB struct {
	dbconn    *gorm.DB
	tableName string
}

//Setup our db by automigrating tables over
func (db *DB) Setup() error {
	conn := db.dbconn

	conn.AutoMigrate(&AdminAccount{})

	if !conn.HasTable(db.tableName) {
		return trace.NotFound("Table not found in DB, Migration fail")
	}
	return nil
}

//SaveInDB saves an admin acc in db
func (db *DB) SaveInDB(a *AdminAccount) error {
	dbconn := db.dbconn
	_, err := db.FindAdmin(a.AccountID)

	if trace.IsNotFound(err) {
		err = dbconn.Create(&a).Error
	} else if err == nil {
		err = db.updateByID(a)
	}

	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (db *DB) UpdateUserCountByID(a *AdminAccount) error {
	conn := db.dbconn

	_, err := db.FindAdmin(a.AccountID)

	if trace.IsNotFound(err) {
		err = conn.Create(&a).Error
	} else {
		err = conn.Table(db.tableName).Where("account_id = ?", a.AccountID).Updates(&AdminAccount{
			Users: a.Users,
		}).Error
	}
	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

//FindAdmin returns the admin user data from the DB
func (db *DB) FindAdmin(AccountID string) (*AdminAccount, error) {
	adminFound := &AdminAccount{}
	conn := db.dbconn
	record := conn.Table(db.tableName).Where("account_id = ?", AccountID).Find(&adminFound)
	if record.RecordNotFound() {
		return nil, trace.NotFound("RecordNotFound")
	}

	// All unahndled errors e.g. db conn errs
	if record.Error != nil {
		return nil, trace.Wrap(record.Error)
	}
	if adminFound.AccountID != AccountID {
		return nil, trace.NotFound("Acc id doenst match")
	}
	return adminFound, nil
}

func (db *DB) updateByID(a *AdminAccount) error {
	conn := db.dbconn
	err := conn.Table(db.tableName).Where("account_id = ?", a.AccountID).Updates(&AdminAccount{
		SessionExpiry: a.SessionExpiry,
		SessionToken:  a.SessionToken,
		IsUpgraded:    a.IsUpgraded,
		MaxUsers:      a.MaxUsers,
	})
	if err.Error != nil {
		return trace.Wrap(err.Error)
	}
	return nil
}

func (db *DB) resetDB() {
	conn := db.dbconn
	conn.DropTableIfExists(&AdminAccount{})
}
