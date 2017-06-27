// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/toomore/mailbox/utils"
)

var (
	conn   *sql.DB
	dryRun *bool
	path   *string
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

func readUser(group string) {
	conn = utils.GetConn()
	rows, err := conn.Query(`SELECT email,groups,created FROM user where groups=?`, group)
	defer rows.Close()
	if err != nil {
		log.Fatal(">>>>>", err)
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

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "user info",
	Long:  `Import user from csv`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("user called")
	},
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import user from csv",
	Long:  "Import user data from csv file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("user.import", args, *path)
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
	},
}

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show users",
	Long:  "Show all/group users",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println(args)
		if len(args) == 0 {
			cmd.Help()
		} else {
			for _, g := range args {
				fmt.Printf("----- %s -----\n", g)
				readUser(g)
			}
		}
	},
}

func init() {
	path = importCmd.Flags().StringP("path", "p", "./list.csv", "csv file path")
	dryRun = importCmd.Flags().BoolP("dryRun", "d", false, "Dry run read csv data")

	RootCmd.AddCommand(userCmd)
	userCmd.AddCommand(importCmd, showCmd)
}
