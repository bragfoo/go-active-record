package main

import (
	"fmt"

	"github.com/bragfoo/go-active-record/db"
)

func main() {
	db.Init("root", "", "127.0.0.1", "mysql", 3306)
	rds, _ := db.Query("select * from user")
	for _, rd := range rds {
		user, _ := rd.Get("user")
		fmt.Println(user)
		fmt.Println(rd)
	}
}
