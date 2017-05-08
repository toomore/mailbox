// mailbox_campaign - create, list, hash campaign info.
/*

Usage:

	mailbox_campaign [flags ...] [cmd]

The cmd are:

	`create`: create a campaign
	`list`: list all campaign info
	`hash`: make a hash with c(cid), u(uid)

The flags are:

	`-c`: cid
	`-u`: uid

*/
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"text/tabwriter"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/toomore/mailbox/campaign"
	"github.com/toomore/mailbox/utils"
)

var (
	conn *sql.DB
	cid  = flag.String("cid", "", "campaign id")
	uid  = flag.String("uid", "", "User id")
)

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

func makeHash() {
	data := url.Values{}
	data.Add("c", *cid)
	data.Add("u", *uid)
	log.Printf("/read/%x?%s\n", campaign.MakeMac(*cid, data), data.Encode())
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
		case "hash":
			makeHash()
		}
	} else {
		fmt.Println("mailbox_campaign [cmd]\ncmd: `create`, `list`, `hash`")
	}
}
