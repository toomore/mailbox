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
	"log"
	"net/url"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/toomore/mailbox/campaign"
	"github.com/toomore/mailbox/utils"
)

var (
	campaignUID *string
	campaignCID *string
)

func create() ([8]byte, [8]byte) {
	id, seed := utils.GenSeed(), utils.GenSeed()
	_, err := conn.Query(fmt.Sprintf(`INSERT INTO campaign(id,seed) VALUES('%s', '%s')`, id, seed))
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
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "id", "seed", "created", "updated*")
	for rows.Next() {
		if err := rows.Scan(&id, &seed, &created, &updated); err != nil {
			log.Println("[err]", err)
		} else {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", id, seed, created, updated)
		}
	}
	w.Flush()
}

func makeHash(cid, uid *string) {
	data := url.Values{}
	data.Add("c", *cid)
	data.Add("u", *uid)
	log.Printf("https://%s/read/%x?%s\n", os.Getenv("mailbox_web_site"), campaign.MakeMac(*cid, data), data.Encode())
}

func openGroups(cid string, groups string) {
	rows, err := conn.Query(`
	SELECT id,email,f_name,reader.created
	FROM user
	LEFT JOIN reader ON (id=reader.uid AND reader.cid=?)
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
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", "id*", "email", "f_name", "open")
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

func openCount(cid string, groups string) {
	rows, err := conn.Query(`
	SELECT uid,u.email,count(*) AS count, min(reader.created) as open, max(reader.created) as latest
	FROM reader, user AS u
	WHERE uid=u.id AND cid=? AND u.groups=?
	GROUP BY uid
	ORDER BY count DESC`, cid, groups)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		count  int
		email  string
		nums   int
		fopen  string
		latest string
		sum    int
		uid    string
	)
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", "uid", "email", "count*", "open", "latest")
	for rows.Next() {
		if err := rows.Scan(&uid, &email, &count, &fopen, &latest); err != nil {
			log.Println("[err]", err)
		} else {
			sum += count
			nums++
			fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n", uid, email, count, fopen, latest)
		}
	}
	fmt.Fprintf(w, "%d\t%.02f%%\t%d\n", nums, float64(sum)/float64(nums)*100, sum)
	w.Flush()
}

// campaignCmd represents the campaign command
var campaignCmd = &cobra.Command{
	Use:   "campaign",
	Short: "campaign operator",
	Long:  `campaign operator`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a campaign",
	Long:  `create a campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		id, seed := create()
		log.Printf("id: %s, seed: %s", id, seed)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list campaign",
	Long:  `list campaign`,
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "hash cid, uid",
	Long:  `hash cid, uid`,
	Run: func(cmd *cobra.Command, args []string) {
		if *campaignCID == "" || *campaignUID == "" {
			cmd.Help()
			log.Fatal("Vars lost `uid`, `cid`")
		}
		makeHash(campaignCID, campaignUID)
	},
}

var openCmd = &cobra.Command{
	Use:   "open [group] [cid ...]",
	Short: "campaign open by group by cid",
	Long:  `campaign open by group by cid`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			log.Fatal("Lost data")
		}
		for _, cid := range args[1:] {
			fmt.Printf("----- %s -----\n", cid)
			openGroups(cid, args[0])
		}
	},
}

var opencountCmd = &cobra.Command{
	Use:   "opencount [group] [cid ...]",
	Short: "count campaign open and list first/latest open by group by cid",
	Long:  `count campaign open and list first/latest open by group by cid`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			log.Fatal("Lost data")
		}
		for _, cid := range args[1:] {
			fmt.Printf("----- %s -----\n", cid)
			openCount(cid, args[0])
		}
	},
}

func init() {
	campaignCID = hashCmd.Flags().String("cid", "", "campaign ID")
	campaignUID = hashCmd.Flags().String("uid", "", "user ID")

	RootCmd.AddCommand(campaignCmd)
	campaignCmd.AddCommand(createCmd, listCmd, hashCmd, openCmd, opencountCmd)
}
