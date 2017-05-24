// mailbox_server - to receive open click.
/*
Usage:

	mailbox_server [flags]

The flags are:

	`-p`: Http port, default is `:8801`

*/
package main

import (
	"encoding/hex"
	"flag"
	"log"
	"net/http"
	"net/url"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
	"github.com/toomore/mailbox/campaign"
	"github.com/toomore/mailbox/reader"
	"github.com/toomore/mailbox/utils"
)

var (
	httpPort = flag.String("p", ":8801", "Http port")
	cExpr    = regexp.MustCompile(`/(read|door)/([0-9a-z]+)`)
)

func read(w http.ResponseWriter, r *http.Request) {
	v, _ := url.ParseQuery(r.Header.Get("X-Args"))

	var hm string
	match := cExpr.FindStringSubmatch(r.Header.Get("X-Uri"))
	if len(match) >= 3 {
		hm = match[2]
	}
	w.WriteHeader(http.StatusNotFound)
	hmbyte, _ := hex.DecodeString(hm)
	if campaign.CheckMac(hmbyte, v.Get("c"), v) {
		log.Println("[read] Pass")
		reader.Save(v.Get("c"), v.Get("u"), r.Header.Get("X-Real-Ip"), r.Header.Get("User-Agent"))
	} else {
		log.Println("[read] Hash Fail!!!")
	}
	log.Printf("%+v", r)
	log.Println(r.Header.Get("X-Real-Ip"))
	log.Println(r.Header.Get("User-Agent"))
}

func door(w http.ResponseWriter, r *http.Request) {
	v, _ := url.ParseQuery(r.Header.Get("X-Args"))
	if v.Get("t") != "a" {
		log.Println("No `t`", v.Get("t"))
		return
	}

	var hm string
	match := cExpr.FindStringSubmatch(r.Header.Get("X-Uri"))
	log.Println(match)
	if len(match) >= 3 {
		hm = match[2]
	}
	hmbyte, _ := hex.DecodeString(hm)
	if campaign.CheckMac(hmbyte, v.Get("c"), v) {
		utils.GetConn().Query(`INSERT INTO doors(cid,uid,linkid,ip,agent) VALUES(?,?,?,?,?)`,
			v.Get("c"), v.Get("u"), v.Get("l"), r.Header.Get("X-Real-Ip"), r.Header.Get("User-Agent"))
		log.Println("[door] Pass")
	} else {
		log.Println("[door] Hash Fail!!!")
	}
	log.Printf("%+v", r)
	log.Println(r.Header.Get("X-Real-Ip"))
	log.Println(r.Header.Get("User-Agent"))
	rows, err := utils.GetConn().Query(`SELECT url FROM links WHERE cid=? AND id=?`, v.Get("c"), v.Get("l"))
	if err == nil {
		for rows.Next() {
			var url string
			rows.Scan(&url)
			http.Redirect(w, r, url, http.StatusSeeOther)
		}
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/read/", read)
	http.HandleFunc("/door/", door)
	log.Println("HTTP Port:", *httpPort)
	log.Println(http.ListenAndServe(*httpPort, nil))
}
