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
	`-rl`: Replace A tag links (default is true)

`-uid`, `-g` can't use together

Example:

	mailbox_sender -cid cbc6eb46 -g testuser -p ./email_1.html -t "#1 New Paper!" -d

*/
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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
	uid         = flag.String("uid", "", "User ID, support more by splited with ','")
)

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
	var rows *sql.Rows
	if *uid != "" {
		uids := strings.Split(*uid, ",")
		for i, v := range uids {
			uids[i] = fmt.Sprintf("'%s'", v)
		}
		rows, err = utils.GetConn().Query(fmt.Sprintf(`SELECT id,email,f_name,l_name FROM user WHERE alive=1 AND id IN (%s)`, strings.Join(uids, ",")))
	} else {
		rows, err = utils.GetConn().Query(`SELECT id,email,f_name,l_name FROM user WHERE alive=1 AND groups=?`, *groups)
	}
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	mails.ProcessSend(body, rows, *cid, *replaceLink, *subject, *dryRun)
}
