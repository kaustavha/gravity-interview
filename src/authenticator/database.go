package authenticator

import (
	"github.com/gravitational/trace"
	"github.com/jinzhu/gorm"
)

type DB struct {
	dbconn *gorm.DB
}

func (db *DB) SaveInDB(a *AdminAccount) {
	dbconn := db.dbconn
	_, err := db.FindAdmin(a.AccountId)

	if trace.IsNotFound(err) {
		dbconn.Create(&a)
	} else if err == nil {
		db.UpdateById(*a)
	}
}

func (db *DB) UpdateById(a AdminAccount) {
	conn := db.dbconn
	conn.Table("admin_accounts").Where("account_id = ?", a.AccountId).Updates(a)
}
func (db *DB) FindAdmin(accountId string) (*AdminAccount, error) {
	adminFound := &AdminAccount{}
	conn := db.dbconn
	record := conn.Table("admin_accounts").Where("account_id = ?", accountId).Find(&adminFound)
	if record.RecordNotFound() {
		return nil, trace.NotFound("RecordNotFound")
	}
	if record.Error != nil {
		return nil, record.Error
	}
	if adminFound.AccountId != accountId {
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
