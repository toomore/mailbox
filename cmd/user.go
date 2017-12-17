// Copyright © 2017 Toomore Chiang
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/toomore/mailbox/utils"
)

var (
	userDryRun *bool
	userConn   *sql.DB
)

type user struct {
	email  string
	groups string
	fname  string
	lname  string
	alive  int
}

func readCSV(path string) []user {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("[cmd][readCSV][open]", err)
	}
	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatal("[cmd][readCSV][ReadAll] ", err)
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
		case "alive":
			for di, dv := range data[1:] {
				alive, err := strconv.Atoi(dv[i])
				if err != nil {
					log.Fatal("Alive type fail", err)
				}
				result[di].alive = alive
			}
		}
	}
	return result
}

func insertInto(data []user) {
	stmt, err := userConn.Prepare(`INSERT INTO user(email,groups,f_name,l_name)
	                           VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE f_name=?, l_name=?`)
	if err != nil {
		log.Fatal("[cmd][insertInto][Prepare]", err)
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

func updateUser(data []user) {
	for _, v := range data {
		rows, err := userConn.Query(`SELECT count(*) AS c FROM user WHERE groups=? AND email=?`, v.groups, v.email)
		if err != nil {
			log.Fatal("[cmd][updateUser][Prepare]", err)
		}
		var c int
		for rows.Next() {
			rows.Scan(&c)
		}
		if c > 0 {
			if result, err := userConn.Exec(`UPDATE user SET f_name=?, l_name=?, alive=? WHERE groups=? AND email=?`,
				v.fname, v.lname, v.alive, v.groups, v.email); err == nil {
				insertID, _ := result.LastInsertId()
				rowAff, _ := result.RowsAffected()
				log.Println("[UPDATE] LastInsertId", insertID, "RowsAffected", rowAff, "email", v.email)
			} else {
				log.Println("[Err]", err)
			}
		} else {
			if v.alive == 1 {
				if result, err := userConn.Exec(`INSERT INTO user(email, groups, f_name, l_name, alive) VALUES(?,?,?,?,?)`,
					v.email, v.groups, v.fname, v.lname, v.alive); err == nil {
					insertID, _ := result.LastInsertId()
					rowAff, _ := result.RowsAffected()
					log.Println("[INSERT] LastInsertId", insertID, "RowsAffected", rowAff, "email", v.email)
				} else {
					log.Println("[Err]", err)
				}
			} else {
				log.Println("[No INSERT alive=0 ]", v.email)
			}
		}
	}
}

func readUser(group string) {
	rows, err := userConn.Query(`SELECT id,email,f_name,l_name,alive,created FROM user WHERE alive=1 AND groups=?`, group)
	defer rows.Close()
	if err != nil {
		log.Fatal("[cmd][readUser][Query]", err)
	}
	var (
		id      string
		email   string
		fname   string
		lname   string
		created time.Time
		alive   int
	)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", "id", "email", "fname", "lname", "alive", "created")
	for rows.Next() {
		if err := rows.Scan(&id, &email, &fname, &lname, &alive, &created); err != nil {
			log.Println(err)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\n", id, email, fname, lname, alive, created)
		}
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", "id", "email", "fname", "lname", "alive", "created")
	w.Flush()
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "User info",
	Long:  `相關 user 的操作`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var importCmd = &cobra.Command{
	Use:   "import [csv path ...]",
	Short: "Import user from csv",
	Long: `匯入使用者資訊，CSV 檔案需要 email, groups, f_name, l_name 欄位， 支援
多檔案匯入。`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		userConn = utils.GetConn()
	},
	Run: func(cmd *cobra.Command, args []string) {
		for n, path := range args {
			log.Printf(">>> Read csv[%d]: `%s`", n, path)
			if *userDryRun {
				log.Println(">>> Dry Run data")
				for i, v := range readCSV(path) {
					fmt.Printf("%d %+v\n", i, v)
				}
			} else {
				insertInto(readCSV(path))
			}
		}
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [csv path ...]",
	Short: "Update user from csv",
	Long: `更新使用者資訊，CSV 檔案需要 email, groups, f_name, l_name, alive 欄位， 支援
多檔案匯入。`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		userConn = utils.GetConn()
	},
	Run: func(cmd *cobra.Command, args []string) {
		for n, path := range args {
			log.Printf(">>> Read csv[%d]: `%s`", n, path)
			if *userDryRun {
				log.Println(">>> Dry Run data")
				for i, v := range readCSV(path) {
					fmt.Printf("%d %+v\n", i, v)
				}
			} else {
				updateUser(readCSV(path))
			}
		}
	},
}

var showCmd = &cobra.Command{
	Use:   "show [groups ...]",
	Short: "Show users",
	Long:  "列出群組使用者名單。",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		userConn = utils.GetConn()
	},
	Run: func(cmd *cobra.Command, args []string) {
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
	userDryRun = userCmd.PersistentFlags().BoolP("dryRun", "d", false, "Dry run read csv data")

	RootCmd.AddCommand(userCmd)
	userCmd.AddCommand(importCmd, showCmd, updateCmd)
}
