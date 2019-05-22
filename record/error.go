package record

import "fmt"

const (
	ErrValueConvert = iota
	ErrColNotFind
	ErrEmptyName
	ErrParseRows
)

type errorActiveRecord struct {
	code     int
	colName  string
	typeName string
}

func (e *errorActiveRecord) Error() string {
	head := "pkg:active_record,msg:"
	switch e.code {
	case ErrParseRows:
		return head + "error to parse rows"
	case ErrEmptyName:
		return head + "string col is empty"
	case ErrValueConvert:
		return head + fmt.Sprintf("col '%s' can not convert type '%s'",
			e.colName, e.typeName)
	case ErrColNotFind:
		return head + fmt.Sprintf("col '%s' not find in rows", e.typeName)
	default:
		return head + "undefined error"
	}
}
