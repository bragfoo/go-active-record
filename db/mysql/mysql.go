package mysql

import (
	"fmt"

	"database/sql"

	"github.com/bragfoo/go-active-record/record"
	_ "github.com/go-sql-driver/mysql"
)

type Mysql struct {
	dbcp *sql.DB
}

func New(user, passwd, host, dbName string, port int) (*Mysql, error) {
	var err error
	db := &Mysql{}
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, passwd,
		host, port, dbName)
	db.dbcp, err = sql.Open("mysql", source)
	if err != nil {
		return nil, err
	}
	err = db.dbcp.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Mysql) Query(sql string, args ...interface{}) (record.ActiveRecordList,
	error) {
	rows, err := db.dbcp.Query(sql, args...)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	rds, err := record.GetActiveRecordList(rows)
	return rds, err
}

func (db *Mysql) QueryFirst(sql string, args ...interface{}) (*record.ActiveRecord,
	error) {
	rows, err := db.dbcp.Query(fmt.Sprintf("%s limit 1", sql), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return record.GetActiveRecord(rows)
}

func (db *Mysql) Update(sql string, args ...interface{}) (int64, error) {
	stmt, err := db.dbcp.Prepare(sql)
	defer stmt.Close()
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

func (db *Mysql) Stat() sql.DBStats {
	return db.dbcp.Stats()
}
func (db *Mysql) Close() error {
	return db.dbcp.Close()
}
