// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/toomore/mailbox/campaign"
	"github.com/toomore/mailbox/mails"
	"github.com/toomore/mailbox/reader"
	"github.com/toomore/mailbox/utils"
)

var (
	servercExpr      = regexp.MustCompile(`/(read|door|washi)/([0-9a-z]+)`)
	serverhttpPort   *string
	serverlinksCache = make(map[string]string)
	serverConn       *sql.DB
)

func serverLog(note string, r *http.Request) {
	log.Printf("%s [%s] \"%+v\" %s\n", note, r.Header.Get("X-Real-Ip"), r, r.Header.Get("User-Agent"))
}

func read(w http.ResponseWriter, r *http.Request) {
	v, _ := url.ParseQuery(r.Header.Get("X-Args"))

	var hm string
	match := servercExpr.FindStringSubmatch(r.Header.Get("X-Uri"))
	if len(match) >= 3 {
		hm = match[2]
	}
	w.WriteHeader(http.StatusNotFound)
	hmbyte, _ := hex.DecodeString(hm)
	if campaign.CheckMac(hmbyte, v.Get("c"), v) {
		serverLog("[read] Pass", r)
		reader.Save(v.Get("c"), v.Get("u"), r.Header.Get("X-Real-Ip"), r.Header.Get("User-Agent"))
	} else {
		serverLog("[read] Hash Fail!!!", r)
	}
}

func washi(v url.Values, url string) []byte {
	washigroup := regexp.MustCompile(`{{WASHI}}(.+){{/WASHI}}`).FindStringSubmatch(url)
	if len(washigroup) > 1 {
		userrows, err := serverConn.Query(`SELECT f_name, l_name FROM user WHERE id=?`, v.Get("u"))
		if err == nil {
			washiURL := []byte(washigroup[1])
			for userrows.Next() {
				var (
					fname string
					lname string
				)
				userrows.Scan(&fname, &lname)
				mails.ReplaceFname(&washiURL, fname)
				mails.ReplaceLname(&washiURL, lname)
				return washiURL
			}
		}
	}
	return nil
}

func door(w http.ResponseWriter, r *http.Request) {
	v, _ := url.ParseQuery(r.Header.Get("X-Args"))

	match := servercExpr.FindStringSubmatch(r.Header.Get("X-Uri"))
	if match[1] == "door" {
		if v.Get("t") != "a" {
			log.Println("No `t`", v.Get("t"))
			return
		}
	}

	var hm string
	if len(match) >= 3 {
		hm = match[2]
	}

	hmbyte, _ := hex.DecodeString(hm)
	if campaign.CheckMac(hmbyte, v.Get("c"), v) {
		serverConn.Query(`INSERT INTO doors(cid,uid,linkid,ip,agent) VALUES(?,?,?,?,?)`,
			v.Get("c"), v.Get("u"), v.Get("l"), r.Header.Get("X-Real-Ip"), r.Header.Get("User-Agent"))
		serverLog("[door] Pass", r)

		var (
			serverlinksCacheKey = fmt.Sprintf("%s|%s", v.Get("c"), v.Get("l"))
			ok                  bool
			url                 string
		)

		if url, ok = serverlinksCache[serverlinksCacheKey]; ok {
			log.Println("Using", match[1], "cache", serverlinksCacheKey, url)
		} else {
			rows, err := serverConn.Query(`SELECT url FROM links WHERE cid=? AND id=?`, v.Get("c"), v.Get("l"))
			if err == nil {
				for rows.Next() {
					rows.Scan(&url)
					serverlinksCache[serverlinksCacheKey] = url
					log.Println("Find", match[1], serverlinksCacheKey, url)
				}
			}
		}

		if url != "" {
			switch match[1] {
			case "door":
				http.Redirect(w, r, url, http.StatusSeeOther)
				return
			case "washi":
				http.Redirect(w, r, string(washi(v, url)), http.StatusSeeOther)
				return
			}
		}
	} else {
		serverLog("[door] Hash Fail!!!", r)
	}
	w.WriteHeader(http.StatusNotFound)
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run tiny server for open, click trace",
	Long:  `啟動一個 web server，來接收開信、點擊連結紀錄。`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		serverConn = utils.GetConn()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run server ...")
		http.HandleFunc("/read/", read)
		http.HandleFunc("/door/", door)
		http.HandleFunc("/washi/", door)
		log.Println("HTTP Port:", *serverhttpPort)
		log.Println(http.ListenAndServe(*serverhttpPort, nil))
	},
}

func init() {
	serverhttpPort = serverCmd.Flags().StringP("port", "p", ":8801", "HTTP Port")
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
