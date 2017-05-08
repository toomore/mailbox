// mailbox_sender - sender worker.
/*
Usage:

	mailbox_sender [flags]

The flags are:

	`-cid`: Campaign id
	`-g`: User groups
	`-p`: HTML file path
	`-t`: Mail Subject
	`-d`: Dry run all but not to send mail

Example:

	mailbox_sender -cid cbc6eb46 -g testuser -p ./email_1.html -t "#1 New Paper!" -d

*/
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
	groups  = flag.String("g", "", "User groups")
	path    = flag.String("p", "", "HTML file path")
	subject = flag.String("t", "", "mail subject")
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
		var (
			email string
			fname string
			lname string
			msg   []byte
			no    string
		)
		rows.Scan(&no, &email, &fname, &lname)

		msg = body
		replaceReader(&msg, *cid, seed, no)
		params := mails.GenParams(
			fmt.Sprintf("%s %s <%s>", fname, lname, email),
			string(msg),
			*subject)
		if !*dryRun {
			mails.Send(params)
		}
		count++
	}
	log.Printf("\n  cid: %s, groups: %s, count: %d\n  Subject: `%s`\n", *cid, *groups, count, *subject)
}
