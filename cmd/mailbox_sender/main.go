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
	`-uid`: User ID

`-uid`, `-g` can't use together

Example:

	mailbox_sender -cid cbc6eb46 -g testuser -p ./email_1.html -t "#1 New Paper!" -d

*/
package main

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/toomore/mailbox/campaign"
	"github.com/toomore/mailbox/mails"
	"github.com/toomore/mailbox/utils"
)

var (
	cid         = flag.String("cid", "", "campaign ID")
	dryRun      = flag.Bool("d", false, "Dry run")
	replaceLink = flag.Bool("rl", true, "Replace A tag links")
	groups      = flag.String("g", "", "User groups")
	path        = flag.String("p", "", "HTML file path")
	subject     = flag.String("t", "", "mail subject")
	uid         = flag.String("uid", "", "User ID")
	areg        = regexp.MustCompile(`href="(http[s]?://[a-zA-z0-9/\.:?=,-]+)"`)
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

func replaceFname(html *[]byte, fname string) {
	*html = bytes.Replace(*html, []byte("{{FNAME}}"), []byte(fname), -1)
}

func replaceATag(html *[]byte, allATags []linksData, cid string, seed string, uid string) {
	for _, v := range allATags {
		data := url.Values{}
		data.Set("c", cid)
		data.Set("u", uid)
		data.Set("l", v.linkID)
		data.Set("t", "a")
		hm := campaign.MakeMacSeed(seed, data)

		*html = bytes.Replace(*html, v.url,
			[]byte(fmt.Sprintf("https://%s/door/%x?%s", os.Getenv("mailbox_web_site"), hm, data.Encode())), -1)
	}
}

type linksData struct {
	md5h   string
	linkID string
	url    []byte
}

func filterATags(body []byte) []linksData {
	allATags := areg.FindAllSubmatch(body, -1)
	result := make([]linksData, len(allATags))
	for i, v := range allATags {
		md5h := md5.New()
		md5h.Write(v[1])
		md5hstr := fmt.Sprintf("%x", md5h.Sum(nil))
		linkID := fmt.Sprintf("%s", utils.GenSeed())
		_, err := utils.GetConn().Query(`INSERT INTO links(id,cid,url,urlhash) VALUES(?,?,?,?)`, linkID, *cid, v[1], md5hstr)
		if err != nil {
			rows, _ := utils.GetConn().Query(`SELECT id FROM links WHERE cid=? AND urlhash=?`, *cid, md5hstr)
			for rows.Next() {
				rows.Scan(&linkID)
			}
		}
		result[i].md5h = md5hstr
		result[i].linkID = linkID
		result[i].url = v[1]
	}
	return result
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
	var rows *sql.Rows
	if *uid != "" {
		rows, err = utils.GetConn().Query(`SELECT id,email,f_name,l_name FROM user WHERE alive=1 AND id=?`, *uid)
	} else {
		rows, err = utils.GetConn().Query(`SELECT id,email,f_name,l_name FROM user WHERE alive=1 AND groups=?`, *groups)
	}
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	var allATags []linksData
	if *replaceLink {
		allATags = filterATags(body)
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
		if *replaceLink {
			replaceATag(&msg, allATags, *cid, seed, *uid)
		}
		replaceFname(&msg, fname)
		replaceReader(&msg, *cid, seed, no)
		params := mails.GenParams(
			fmt.Sprintf("%s %s <%s>", fname, lname, email),
			string(msg),
			*subject)
		if *dryRun {
			log.Printf("%s\n", msg)
			for i, v := range allATags {
				fmt.Printf("%d => [%s] %s\n", i, v.linkID, v.url)
			}
		} else {
			mails.Send(params)
		}
		count++
	}
	if *uid != "" {
		log.Printf("\n  cid: %s, uid: %s, count: %d\n  Subject: `%s`\n", *cid, *uid, count, *subject)
	} else {
		log.Printf("\n  cid: %s, groups: %s, count: %d\n  Subject: `%s`\n", *cid, *groups, count, *subject)
	}
}
