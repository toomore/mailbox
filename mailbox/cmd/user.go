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
	Use:   "import [csv path ...]",
	Short: "Import user from csv",
	Long:  "Import user data from csv file",
	Run: func(cmd *cobra.Command, args []string) {
		for n, path := range args {
			log.Printf(">>> Read csv[%d]: `%s`", n, path)
			if *dryRun {
				log.Println(">>> Dry Run data")
				for i, v := range readCSV(path) {
					fmt.Printf("%d %+v\n", i, v)
				}
			} else {
				conn = utils.GetConn()
				insertInto(readCSV(path))
			}
		}
	},
}

var showCmd = &cobra.Command{
	Use:   "show [groups ...]",
	Short: "Show users",
	Long:  "Show all/group users",
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
	dryRun = importCmd.Flags().BoolP("dryRun", "d", false, "Dry run read csv data")

	RootCmd.AddCommand(userCmd)
	userCmd.AddCommand(importCmd, showCmd)
}
