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
	var (
		id      string
		seed    string
		created time.Time
		updated time.Time
	)
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

func openGroups(cid string, groups string) {
	rows, err := conn.Query(`
	SELECT id,email,f_name,reader.created
	FROM user
	LEFT JOIN reader ON (id=reader.uid and reader.cid=?)
	WHERE groups=?
	GROUP BY id;`, cid, groups)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		id        string
		email     string
		fname     string
		created   sql.NullString
		nums      int
		openCount int
	)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "id", "email", "f_name", "open")
	for rows.Next() {
		if err := rows.Scan(&id, &email, &fname, &created); err != nil {
			log.Println("[err]", err)
		} else {
			nums++
			if created.Valid {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", id, email, fname, created.String)
				openCount++
			} else {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", id, email, fname, "Not open")
			}
		}
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t%.2f%%\n", "", "", "", float64(openCount)/float64(nums)*100)
	w.Flush()
}

func openList(cid string) {
	rows, err := conn.Query(`
	SELECT uid,u.email,count(*) AS count
	FROM reader, user AS u
	WHERE uid=u.id AND cid=?
	GROUP BY uid
	ORDER BY count DESC`, cid)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		count int
		email string
		nums  int
		sum   int
		uid   string
	)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "%s\t%s\t%s\n", "uid", "email", "count")
	for rows.Next() {
		if err := rows.Scan(&uid, &email, &count); err != nil {
			log.Println("[err]", err)
		} else {
			sum += count
			nums++
			fmt.Fprintf(w, "%s\t%s\t%d\n", uid, email, count)
		}
	}
	fmt.Fprintf(w, "%d\t%.02f%%\t%d\n", nums, float64(sum)/float64(nums)*100, sum)
	w.Flush()
}

func printTips() {
	fmt.Println("mailbox_campaign [cmd]\ncmd: `create`, `list`, `hash`, `open [cid] [groups]`, `openlist [cid]`")
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
		case "open":
			if len(args) >= 3 {
				openGroups(args[1], args[2])
			} else {
				printTips()
			}
		case "openlist":
			if len(args) >= 2 {
				openList(args[1])
			} else {
				printTips()
			}
		default:
			printTips()
		}
	} else {
		printTips()
	}
}
