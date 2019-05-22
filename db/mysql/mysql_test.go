package mysql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestNew(t *testing.T) {
	type args struct {
		user   string
		passwd string
		host   string
		db     string
		port   int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"A",
			args{
				"root",
				"",
				"localhost",
				"test",
				3306,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := New(tt.args.user, tt.args.passwd, tt.args.host,
				tt.args.db, tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Init() got Close = %v, error = %v, wantErr %v", err, got.Close(), tt.wantErr)
			}
		})
	}
}

func TestDBMysql_Find(t *testing.T) {
	db, _ := New("root", "", "127.0.0.1", "test", 3306)
	type args struct {
		sql  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"A",
			args{
				sql:  "select * from pools",
				args: []interface{}{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.Query(tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for _, rd := range got {
				t.Logf("%s\n", rd)
			}
		})
	}
}

func TestDBMysql_FindFirst(t *testing.T) {
	db, _ := New("root", "", "127.0.0.1", "test", 3306)
	type args struct {
		sql  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"A",
			args{
				sql:  "select * from pools",
				args: []interface{}{},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.QueryFirst(tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.FindFirst() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", got)
		})
	}
}
