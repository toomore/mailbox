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
	"io"
	"log"
	"os"
	"strconv"
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
	sendTextPath    *string
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
			log.Fatal("[cmd][send][open] ", err)
		}

		file_text, err := os.Open(*sendTextPath)
		if err != nil {
			log.Fatal("[cmd][send][open] ", err)
		}

		body, err := io.ReadAll(file)
		if err != nil {
			log.Fatal("[cmd][send][ReadAll] ", err)
		}
		body_text, err := io.ReadAll(file_text)
		if err != nil {
			log.Fatal("[cmd][send][ReadAll] ", err)
		}

		var rows *sql.Rows
		if *sendUID != "" {
			uids := strings.Split(*sendUID, ",")
			placeholders := make([]string, len(uids))
			args := make([]interface{}, len(uids))
			for i, v := range uids {
				uid := strings.TrimSpace(v)
				if _, err := strconv.Atoi(uid); err != nil {
					log.Fatal("[cmd][send] invalid uid: ", v)
				}
				placeholders[i] = "?"
				args[i] = uid
			}
			query := `SELECT id,email,f_name,l_name FROM user WHERE alive=1 AND id IN (` + strings.Join(placeholders, ",") + `)`
			rows, err = sendConn.Query(query, args...)
		} else {
			rows, err = sendConn.Query(`SELECT id,email,f_name,l_name FROM user WHERE alive=1 AND groups=?`, *sendGroups)
		}
		if err != nil {
			log.Fatal("[cmd][send][Query] ", err)
		}
		if rows != nil {
			defer rows.Close()
		}

		mails.ProcessSend(body, body_text, rows, *sendCID, *sendReplaceLink, *sendSubject, *sendDryRun, *sendLimit)
	},
}

func init() {
	sendCID = sendCmd.Flags().String("cid", "", "Campaign ID")
	sendUID = sendCmd.Flags().String("uid", "", "User ID, support more by splited with ','")
	sendDryRun = sendCmd.Flags().BoolP("dryrun", "d", false, "Dry run")
	sendGroups = sendCmd.Flags().StringP("groups", "g", "", "User groups")
	sendPath = sendCmd.Flags().StringP("path", "p", "", "HTML file path")
	sendTextPath = sendCmd.Flags().StringP("text", "t", "", "Plain file path")
	sendReplaceLink = sendCmd.Flags().Bool("rl", true, "Replace A tag links")
	sendSubject = sendCmd.Flags().StringP("subject", "s", "", "Mail subject")
	sendLimit = sendCmd.Flags().IntP("limit", "", 7, "Send concurrency limit")

	RootCmd.AddCommand(sendCmd)
}
