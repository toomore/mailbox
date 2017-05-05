package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/toomore/mailbox/utils"
)

var conn *sql.DB

func create() ([8]byte, [8]byte) {
	id, seed := utils.GenSeed(), utils.GenSeed()
	_, err := conn.Query(fmt.Sprintf(`INSERT INTO campaign(id,seed) VALUES('%s', '%s')`, id, seed))
	defer conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	return id, seed
}

func list() {
	rows, err := conn.Query(`SELECT id,seed,created,updated FROM campaign ORDER BY updated DESC`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var id string
	var seed string
	var created time.Time
	var updated time.Time
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "id", "seed", "created", "updated")
	for rows.Next() {
		if err := rows.Scan(&id, &seed, &created, &updated); err != nil {
			log.Println("[err]", err)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", id, seed, created, updated)
		}
	}
	w.Flush()
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) >= 1 {
		conn = utils.GetConn()
		switch args[0] {
		case "create":
			id, seed := create()
			log.Printf("id: %s, seed: %s", id, seed)
		case "list":
			list()
		}
	} else {
		fmt.Println("mailbox_campaign [cmd]\ncmd: `create`, `list`")
	}
}
