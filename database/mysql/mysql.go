package mysql

import (
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/bragfoo/go-active-record/database/record"
)

type DB struct {
	dbcp *sql.DB
	once sync.Once
}

func (db *DB) Init(user, passwd, host, dbName string, port int) error {
	var err error
	db.once.Do(func() {
		db.dbcp, err = db.NewDB(user, passwd, host, dbName, port)
	})
	if err != nil {
		return err
	}
	return db.dbcp.Ping()
}

func (db *DB) NewDB(user, passwd, host, dbName string, port int) (*sql.DB,
	error) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd,
		host, port, dbName)
	var err error
	cp, err := sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}
	return cp, nil
}

func (db *DB) Find(sql string, args ...interface{}) (record.ActiveRecordList,
	error) {
	rows, err := db.dbcp.Query(sql, args...)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	rds, err := record.GetActiveRecordList(rows)
	return rds, err
}

func (db *DB) FindFirst(sql string, args ...interface{}) (*record.ActiveRecord,
	error) {
	rows, err := db.dbcp.Query(fmt.Sprintf("%s limit 1", sql), args...)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return record.GetActiveRecord(rows)
}

func (db *DB) Update(sql string, args ...interface{}) (int64, error) {
	stmt, err := db.dbcp.Prepare(sql)
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(args...)
	if err != nil {
		return -1, err
	}
	affect, err := res.RowsAffected()

	if err != nil {
		return -1, err
	}
	return affect, nil
}

func (db *DB) Save(table string, rd *record.ActiveRecord) (bool, error) {
	return false, nil
}

func (db *DB) Delete(table string, rd *record.ActiveRecord) (int64, error) {
	return -1, nil
}

func (db *DB) Close() (error) {
	return db.dbcp.Close()
}
