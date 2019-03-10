package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	dbhost           = "localhost"
	dbport           = "5432"
	dbuser           = "postgres"
	dbname           = "iotdb"
	dbpass           = "mysecretpassword"
	dbsslmode        = "disable"
	defaultTableName = "metrics"
)

type DB struct {
	db *gorm.DB
}

var db *DB

func createDBConn() (*DB, error) {
	db = &DB{db: nil}
	conn, err := gorm.Open("postgres",
		"host="+dbhost+" "+
			"port="+dbport+" "+
			"user="+dbuser+" "+
			"dbname="+dbname+" "+
			"password="+dbpass+" "+
			"sslmode="+dbsslmode)
	if err != nil {
		fmt.Println(err, "db conn err")
		return db, err
	}

	conn.AutoMigrate(&Metric{}, &AdminAccount{})
	if !conn.HasTable(defaultTableName) {
		fmt.Println(conn.HasTable(defaultTableName))
		fmt.Println("Migration fail")
	}

	db = &DB{db: conn}
	return db, err
}
func GetDB() *DB {
	return db
}
func GetDBConn() *gorm.DB {
	return db.db
}

func (db *DB) getConn() *gorm.DB {
	return db.db
}

func (db *DB) countAllUniqueUsersInAccount(accountId string) int {
	return db._countAllUniqueUsersInAccount(defaultTableName, accountId)
}

func (db *DB) _countAllUniqueUsersInAccount(tableName string, accountId string) int {
	conn := db.getConn()
	count := 0
	conn.Table(tableName).Where("account_id = ?", accountId).Count(&count)
	return count
}

func (db *DB) countAllInTable() int {
	return db._countAllInTable(defaultTableName)
}

func (db *DB) _countAllInTable(tableName string) int {
	conn := db.getConn()
	count := 0
	conn.Table(tableName).Count(&count)
	return count
}

func (db *DB) countAll() {

}
