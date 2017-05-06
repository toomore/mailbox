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
)

var httpPort = flag.String("p", ":8801", "Http port")
var readExpr = regexp.MustCompile(`/read/([0-9a-z]+)`)

func read(w http.ResponseWriter, r *http.Request) {
	//u, _ := url.Parse(r.RequestURI)
	//v, _ := url.ParseQuery(u.RawQuery)
	v, _ := url.ParseQuery(r.Header.Get("X-Args"))

	var hm string
	match := readExpr.FindStringSubmatch(r.Header.Get("X-Uri"))
	if len(match) > 1 {
		hm = match[1]
	}
	w.WriteHeader(http.StatusNotFound)
	hmbyte, _ := hex.DecodeString(hm)
	if campaign.CheckMac(hmbyte, v.Get("cid"), v) {
		log.Println("Pass")
	} else {
		log.Println("Hash Fail!!!")
	}
	log.Printf("%+v", r)
	log.Println(r.Header.Get("X-Real-Ip"))
	log.Println(r.Header.Get("User-Agent"))
}

func main() {
	flag.Parse()
	http.HandleFunc("/read/", read)
	log.Println("HTTP Port:", *httpPort)
	log.Println(http.ListenAndServe(*httpPort, nil))
}
