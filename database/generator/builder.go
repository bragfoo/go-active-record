package main

import (
	"fmt"
	"flag"

	_ "github.com/go-sql-driver/mysql"
	"text/template"
	"bufio"
	"os"
	"strings"
	"github.com/bragfoo/go-active-record/database/mysql"
	"github.com/bragfoo/go-active-record/database/record"
)

type col struct {
	Table      string
	Filed      string
	ColName    string
	ColType    string
	ReturnType string
}
type dbData struct {
	Pkg    string
	BTable string
	STable string
	Cols   []*col
}

var dbTemplate = `package {{.Pkg}}

import (
	"github.com/bragfoo/go-active-record/database/record"
	"github.com/bragfoo/go-active-record/database/mysql"
)

type {{.BTable}}RdList []*{{.BTable}}Rd

type {{.BTable}}Rd struct {
	record.ActiveRecord
}

type {{.BTable}} struct {
	db *mysql.DB
}

func (t *{{.BTable}}) Find(sql string, args ...interface{}) ({{.BTable}}RdList,
	error) {
	rds, err := t.db.Find(sql, args...)
	{{.STable}}Rds := make([]*{{.BTable}}Rd, len(rds))
	for i, rd := range rds {
		{{.STable}}Rds[i] = &{{.BTable}}Rd{*rd}
	}
	return {{.STable}}Rds, err
}

func (t *{{.BTable}}) FindFirst(sql string, args ...interface{}) (*{{.BTable}}Rd,
	error) {
	rd, err := t.db.FindFirst(sql, args...)
	return &{{.BTable}}Rd{*rd}, err

}

func (t *{{.BTable}}) Update(sql string, args ...interface{}) (int64, error) {
	return t.db.Update(sql, args)
}

func (t *{{.BTable}}) Save(rd *{{.BTable}}Rd) (bool, error) {
	return t.db.Save("{{.STable}}", &rd.ActiveRecord)
}

func (t *{{.BTable}}) Delete(rd *{{.BTable}}Rd) (int64, error) {
	return t.db.Delete("{{.STable}}", &rd.ActiveRecord)
}

{{range .Cols}}
func (rd *{{.Table}}Rd) Get{{.ColName}}() ({{.ReturnType}}, error) {
	return rd.Get{{.ColType}}("{{.Filed}}")
}
{{end}}
`

func main() {
	var (
		targetPackage string
		dbHost        string
		dbPort        int
		dbName        string
		dbUser        string
		dbPwd         string
	)

	flag.StringVar(&targetPackage, "target_package",
		"github.com/bragfoo/go-active-record/target",
		"target package")
	flag.StringVar(&dbHost, "db_host", "localhost", "db server ip or hostname")
	flag.IntVar(&dbPort, "db_port", 3306, "db port")
	flag.StringVar(&dbName, "db_name", "confman", "db name")
	flag.StringVar(&dbUser, "db_user", "root", "db user")
	flag.StringVar(&dbPwd, "db_pwd", "", "db password")
	flag.Parse()

	db := new(mysql.DB)
	err := db.Init(dbUser, dbPwd, dbHost, dbName, dbPort)
	fck(err)

	tableRds, err := db.Find("show tables")
	fck(err)

	path := getPath(targetPackage)

	for _, tableRd := range tableRds {
		table, err := tableRd.GetString(fmt.Sprintf("Tables_in_%s", dbName))
		fck(err)

		filePath := path + "/" + table + ".go"
		fmt.Printf("gen %s and write file to %s", table, filePath)

		colRds, err := db.Find(fmt.Sprintf("desc %s", table))
		fck(err)

		bTableName := formatName(table, true)
		sTableName := formatName(table, false)
		cols := parseCol(colRds, table)

		w := newWriter(filePath)
		pkgs := strings.Split(targetPackage, "/")
		pkg := pkgs[len(pkgs)-1]

		tmpl, err := template.New("db").Parse(dbTemplate)
		tmpl.Execute(w, &dbData{
			pkg,
			bTableName,
			sTableName,
			cols,
		})

		w.Flush()

	}
}

func newWriter(path string) (*bufio.Writer) {
	f, err := os.Create(path)
	fck(err)
	fck(err)
	w := bufio.NewWriter(f)
	return w
}

func getPath(targetPackage string) string {
	goPath := os.Getenv("GOPATH")
	path := goPath + "/src/" + targetPackage
	if goPath == "" {
		goPath = "~/go"
	}
	os.MkdirAll(path, 0755)
	return path
}

func formatName(str string, cfl bool) string {
	str = strings.ToLower(str)
	if strings.Contains(str, "_") {
		for i, v := range strings.Split(str, "_") {
			if i == 0 {
				str = v
			} else if len(v) > 0 {
				str += strings.ToUpper(v[:1]) + v[1:]
			}
		}
	}

	if cfl {
		return strings.ToUpper(str[:1]) + str[1:]
	}
	return str
}

func formatType(str string) string {
	if strings.Contains(str, "int") {
		return "Int64"
	} else if strings.Contains(str, "varchar") {
		return "String"
	} else {
		return "String"
	}
}

func parseCol(colRds record.ActiveRecordList, table string) []*col {
	cols := make([]*col, len(colRds))
	for i, colRd := range colRds {
		colName, err := colRd.GetString("field")
		fck(err)
		colType, err := colRd.GetString("type")
		cols[i] = &col{
			formatName(table, true),
			colName,
			formatName(colName, true),
			formatType(colType),
			strings.ToLower(formatType(colType)),
		}
	}
	return cols
}

func fck(err error) {
	if err != nil {
		panic(err)
	}
}
