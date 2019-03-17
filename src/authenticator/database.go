package authenticator

import (
	"github.com/gravitational/trace"
	"github.com/jinzhu/gorm"
)

type DB struct {
	dbconn *gorm.DB
}

func (db *DB) SaveInDB(a *AdminAccount) error {
	dbconn := db.dbconn
	_, err := db.FindAdmin(a.AccountID)

	if trace.IsNotFound(err) {
		err = dbconn.Create(&a).Error
	} else if err == nil {
		err = db.updateById(*a)
	}

	if err != nil {
		return trace.Wrap(err)
	}
	return nil
}

func (db *DB) updateById(a AdminAccount) error {
	conn := db.dbconn
	err := conn.Table("admin_accounts").Where("account_id = ?", a.AccountID).Updates(a)
	if err.Error != nil {
		return trace.Wrap(err.Error)
	}
	return nil
}
func (db *DB) FindAdmin(accountId string) (*AdminAccount, error) {
	adminFound := &AdminAccount{}
	conn := db.dbconn
	record := conn.Table("admin_accounts").Where("account_id = ?", accountId).Find(&adminFound)
	if record.RecordNotFound() {
		return nil, trace.NotFound("RecordNotFound")
	}

	// All unahndled errors e.g. db conn errs
	if record.Error != nil {
		return nil, trace.Wrap(record.Error)
	}
	if adminFound.AccountID != accountId {
		return nil, trace.NotFound("Acc id doenst match")
	}
	return adminFound, nil
}

func (db *DB) Setup(defaultTableName string) error {
	conn := db.dbconn
	conn.AutoMigrate(&AdminAccount{})

	if !conn.HasTable(defaultTableName) {
		return trace.NotFound("Table not found in DB, Migration fail")
	}
	return nil
}
