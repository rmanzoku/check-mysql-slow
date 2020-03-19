package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mackerelio/checkers"
)

var (
	hostname        = "127.0.0.1"
	port            = "3306"
	username        = "monitor"
	password        = "monitor"
	time     uint64 = 1
	dsn             = ""
)

const dsnFormat = "%s:%s@tcp(%s:%s)/information_schema"

func run() (err error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	q := "SELECT STATE, TIME, INFO FROM PROCESSLIST WHERE COMMAND <> \"Sleep\" AND TIME = ?"
	rows, err := db.Query(q, time)
	if err != nil {
		return
	}
	defer rows.Close()

	var state string
	var info string
	var t int64
	for rows.Next() {
		err := rows.Scan(&state, &t, &info)
		if err != nil {
			return err
		}
		return fmt.Errorf("Slow query time: %d, info: \"%s\"", t, info)
	}

	err = rows.Err()
	if err != nil {
		return
	}
	return nil
}

func main() {
	flag.StringVar(&hostname, "hostname", hostname, "hostname")
	flag.StringVar(&username, "username", username, "username")
	flag.StringVar(&password, "password", password, "password")
	flag.Uint64Var(&time, "time", time, "time threshold")
	flag.Parse()

	dsn = fmt.Sprintf(dsnFormat, username, password, hostname, port)

	chr := checkers.Ok("command ok!")
	err := run()
	if err != nil {
		chr = checkers.Critical(err.Error())
	}
	chr.Exit()
}
