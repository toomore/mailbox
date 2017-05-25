// mailbox_import_csv - import user data from csv file.
/*

CSV Format:

	email, groups, f_name, l_name

and default import csv path is `./list.csv`

Usage:

	mailbox_import_csv [flags]

The flags are:

	`-p`: csv file path, default is `./list.csv`
	`-d`: using dry run to review csv data

*/
package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/toomore/mailbox/utils"
)

var (
	conn   *sql.DB
	path   = flag.String("p", "./list.csv", "user's csv")
	dryRun = flag.Bool("d", false, "Dry run read csv data")
)

type user struct {
	email  string
	groups string
	fname  string
	lname  string
}

func readCSV(path string) []user {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	result := make([]user, len(data[1:]))
	for i, v := range data[0] {
		switch v {
		case "email":
			for di, dv := range data[1:] {
				result[di].email = dv[i]
			}
		case "groups":
			for di, dv := range data[1:] {
				result[di].groups = dv[i]
			}
		case "f_name":
			for di, dv := range data[1:] {
				result[di].fname = dv[i]
			}
		case "l_name":
			for di, dv := range data[1:] {
				result[di].lname = dv[i]
			}
		}
	}
	return result
}

func insertInto(data []user) {
	stmt, err := conn.Prepare(`INSERT INTO user(email,groups,f_name,l_name)
	                           VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE f_name=?, l_name=?`)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range data {
		if result, err := stmt.Exec(v.email, v.groups, v.fname, v.lname, v.fname, v.lname); err == nil {
			insertID, _ := result.LastInsertId()
			rowAff, _ := result.RowsAffected()
			log.Println("LastInsertId", insertID, "RowsAffected", rowAff)
		} else {
			log.Println("[Err]", err)
		}
	}
}

func readUser() {
	rows, err := conn.Query(`SELECT email,groups,created FROM user;`)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var email string
	var groups string
	var created time.Time
	for rows.Next() {
		if err := rows.Scan(&email, &groups, &created); err != nil {
			log.Println(err)
		}
		log.Println(email, groups, created)
	}
}

func main() {
	flag.Parse()
	log.Printf(">>> Read csv: `%s`", *path)
	data := readCSV(*path)
	if *dryRun {
		log.Println(">>> Dry Run data")
		for i, v := range data {
			fmt.Printf("%d %+v\n", i, v)
		}
	} else {
		conn = utils.GetConn()
		insertInto(data)
	}
	//readUser()
}
