package db

import (
	"sync"

	"github.com/bragfoo/go-active-record/db/mysql"
	"github.com/bragfoo/go-active-record/record"
	"github.com/pkg/errors"
)

// DB interface is portal to use active-record
type DB interface {
	Query(sql string, args ...interface{}) (record.ActiveRecordList, error)
	QueryFirst(sql string, args ...interface{}) (*record.ActiveRecord, error)
	Update(sql string, args ...interface{}) (int64, error)
}

var instance DB
var once sync.Once

func Init(user, passwd, host, dbName string, port int) error {
	var err error
	once.Do(func() {
		instance, err = mysql.New(user, passwd, host, dbName, port)
	})
	if err != nil {
		return errors.Wrap(err, "database connect fails")
	}
	return nil
}

func GetDB() (DB, error) {
	if instance == nil {
		return nil, errors.New("database connection pool is uninitialized")
	}
	return instance, nil
}

func Query(sql string, args ...interface{}) (record.ActiveRecordList,
	error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	return db.Query(sql, args...)
}

func QueryFirst(sql string, args ...interface{}) (*record.ActiveRecord,
	error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}
	return db.QueryFirst(sql, args...)
}

func Update(sql string, args ...interface{}) (int64, error) {
	db, err := GetDB()
	if err != nil {
		return -1, err
	}
	return db.Update(sql, args...)
}
