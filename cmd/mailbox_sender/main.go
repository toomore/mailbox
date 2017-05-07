package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/toomore/mailbox/campaign"
	"github.com/toomore/mailbox/mails"
	"github.com/toomore/mailbox/utils"
)

var (
	cid     = flag.String("cid", "", "campaign ID")
	dryRun  = flag.Bool("d", false, "Dry run")
	groups  = flag.String("g", "", "Groups")
	path    = flag.String("p", "", "HTML file path")
	subject = flag.String("t", "", "Subject")
)

func replaceReader(html *[]byte, cid string, seed string, uid string) {
	data := url.Values{}
	data.Set("c", cid)
	data.Set("u", uid)
	hm := campaign.MakeMacSeed(seed, data)
	*html = bytes.Replace(
		*html,
		[]byte("{{READER}}"),
		[]byte(fmt.Sprintf("https://%s/read/%x?%s", os.Getenv("mailbox_web_site"), hm, data.Encode())),
		1)
}

func main() {
	flag.Parse()
	file, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	seed := campaign.GetSeed(*cid)
	conn := utils.GetConn()
	rows, err := conn.Query(`SELECT id,email,f_name,l_name from user where alive=1 and groups=?`, *groups)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	var count int
	for rows.Next() {
		var no string
		var email string
		var fname string
		var lname string
		rows.Scan(&no, &email, &fname, &lname)

		replaceReader(&body, *cid, seed, no)
		params := mails.GenParams(
			fmt.Sprintf("%s %s <%s>", fname, lname, email),
			string(body),
			*subject)
		if !*dryRun {
			mails.Send(params)
		}
		count++
	}
	log.Printf("\n  cid: %s, groups: %s, count: %d\n  Subject: `%s`\n", *cid, *groups, count, *subject)
}
