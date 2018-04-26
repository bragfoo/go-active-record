package interfaces

import (
	"github.com/bragfoo/go-active-record/database/record"
)

type DB interface {
	Find(sql string, args ...interface{}) ([]*record.ActiveRecord, error)
	FindFirst(sql string, args ...interface{}) (*record.ActiveRecord, error)

	Update(sql string, args ...interface{}) (int64, error)

	Save(table string, rd *record.ActiveRecord) (bool, error)

	Delete(table string, rd *record.ActiveRecord) (int64, error)
}
