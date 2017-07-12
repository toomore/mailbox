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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/toomore/mailbox/mails"
	"github.com/toomore/mailbox/utils"
)

var (
	sendCID         *string
	sendDryRun      *bool
	sendGroups      *string
	sendLimit       *int
	sendPath        *string
	sendReplaceLink *bool
	sendSubject     *string
	sendUID         *string
	sendConn        *sql.DB
)

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send paper",
	Long:  `寄送電子報，處理開信連結與替換點擊連結。`,
	Example: `群組寄送：mailbox send -p {path} -s 'Title: #1' -g {group} --cid {cid}
個人寄送：mailbox send -p {path} -s 'Title: #1' --uid='6,12' --cid {cid}`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		sendConn = utils.GetConn()
	},
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open(*sendPath)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		var rows *sql.Rows
		if *sendUID != "" {
			uids := strings.Split(*sendUID, ",")
			for i, v := range uids {
				uids[i] = fmt.Sprintf("'%s'", v)
			}
			rows, err = sendConn.Query(fmt.Sprintf(`SELECT id,email,f_name,l_name FROM user WHERE alive=1 AND id IN (%s)`, strings.Join(uids, ",")))
		} else {
			rows, err = sendConn.Query(`SELECT id,email,f_name,l_name FROM user WHERE alive=1 AND groups=?`, *sendGroups)
		}
		defer rows.Close()
		if err != nil {
			log.Fatal(err)
		}

		mails.ProcessSend(body, rows, *sendCID, *sendReplaceLink, *sendSubject, *sendDryRun, *sendLimit)
	},
}

func init() {
	sendCID = sendCmd.Flags().String("cid", "", "Campaign ID")
	sendUID = sendCmd.Flags().String("uid", "", "User ID, support more by splited with ','")
	sendDryRun = sendCmd.Flags().BoolP("dryrun", "d", false, "Dry run")
	sendGroups = sendCmd.Flags().StringP("groups", "g", "", "User groups")
	sendPath = sendCmd.Flags().StringP("path", "p", "", "HTML file path")
	sendReplaceLink = sendCmd.Flags().Bool("rl", true, "Replace A tag links")
	sendSubject = sendCmd.Flags().StringP("subject", "s", "", "Mail subject")
	sendLimit = sendCmd.Flags().IntP("limit", "", 7, "Send concurrency limit")

	RootCmd.AddCommand(sendCmd)
}
