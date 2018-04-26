package record

import (
	"database/sql"
	"strconv"
	"fmt"
	"io"
	"bytes"
	"strings"
)

type ActiveRecordList []*ActiveRecord

func (rds ActiveRecordList) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "[")
			i := 0

			for _, rd := range rds {
				fmt.Fprintf(&buf, "{")
				j := 0
				for k, v := range rd.rows {
					fmt.Fprintf(&buf, "\n  %s:%s", k, string((*v)[:]))
					if j++; j < len(rd.rows) {
						fmt.Fprintf(&buf, ",")
					}
				}
				fmt.Fprintf(&buf, "\n}")
				if i++; i < len(rds) {
					fmt.Fprintf(&buf, ",")
				}
			}

			io.WriteString(s, buf.String())
			return
		}

		fallthrough
	case 's', 'q':
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "[")
		i := 0
		for _, rd := range rds {
			fmt.Fprintf(&buf, "{")
			j := 0
			for k, v := range rd.rows {
				fmt.Fprintf(&buf, "%s:%s", k, string((*v)[:]))
				if j++; j < len(rd.rows) {
					fmt.Fprintf(&buf, ",")
				}
			}
			fmt.Fprintf(&buf, "}")
			if i++; i < len(rds) {
				fmt.Fprintf(&buf, ",")
			}
		}
		fmt.Fprintf(&buf, "]")
		io.WriteString(s, buf.String())
	}
}

type ActiveRecord struct {
	table *[]string
	rows  map[string]*sql.RawBytes
}

func GetActiveRecordList(rows *sql.Rows) (ActiveRecordList, error) {
	rds := make([]*ActiveRecord, 0)
	if rows == nil {
		return []*ActiveRecord{}, &errorActiveRecord{ErrParseRows, "", ""}
	}

	for rows.Next() {
		rd, err := parse(rows)
		if err != nil {
			return nil, err
		}
		rds = append(rds, rd)
	}
	return rds, nil
}

func GetActiveRecord(rows *sql.Rows) (*ActiveRecord, error) {
	if rows.Next() {
		return parse(rows)
	}
	return &ActiveRecord{&[]string{}, map[string]*sql.RawBytes{}}, nil
}

func parse(rows *sql.Rows) (*ActiveRecord, error) {
	table, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	lenCN := len(table)
	raw := newRaw(lenCN)
	rows.Scan(raw...)
	rd := &ActiveRecord{
		table: &table,
		rows:  make(map[string]*sql.RawBytes),
	}
	for i := 0; i < lenCN; i++ {
		rd.rows[strings.ToLower(table[i])] = raw[i].(*sql.RawBytes)
	}
	return rd, nil
}

func newRaw(lenCN int) []interface{} {
	raw := make([]interface{}, lenCN)
	for i := 0; i < lenCN; i++ {
		raw[i] = new(sql.RawBytes)
	}
	return raw
}

func (rd *ActiveRecord) Get(colName string) (string,
	error) {
	if colName == "" {
		return "", &errorActiveRecord{ErrEmptyName, "", ""}
	}
	value, ok := rd.rows[strings.ToLower(colName)]
	if !ok {
		return "", &errorActiveRecord{ErrColNotFind, colName, ""}
	} else {
		result := string((*value)[:])
		return result, nil
	}
}

func (rd *ActiveRecord) GetString(colName string) (string, error) {
	return rd.Get(colName)
}

func (rd *ActiveRecord) GetInt(colName string) (int, error) {
	value, err := rd.Get(colName)
	if err != nil {
		return 0, err
	}
	v, err := strconv.Atoi(value)
	if err != nil {
		return 0, &errorActiveRecord{ErrValueConvert, colName, "int"}
	}
	return v, nil
}

func (rd *ActiveRecord) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			var buf bytes.Buffer
			for k, v := range rd.rows {
				fmt.Fprintf(&buf, "%s:  %s,\n", k, string((*v)[:]))
			}
			io.WriteString(s, buf.String())
			return
		}

		fallthrough
	case 's', 'q':
		var buf bytes.Buffer
		for k, v := range rd.rows {
			fmt.Fprintf(&buf, "%s:%s,", k, string((*v)[:]))
		}
		io.WriteString(s, buf.String())
	}
}
