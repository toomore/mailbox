package mails

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/url"
	"os"
	"regexp"

	"github.com/toomore/mailbox/campaign"
	"github.com/toomore/mailbox/utils"
)

var areg = regexp.MustCompile(`href="(http[s]?://[a-zA-z0-9/\.:?=,-@]+)"`)

// ReplaceReader is to replace reader open mail link
func ReplaceReader(html *[]byte, cid string, seed string, uid string) {
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

// ReplaceFname is to replace FNAME tag
func ReplaceFname(html *[]byte, fname string) {
	*html = bytes.Replace(*html, []byte("{{FNAME}}"), []byte(fname), -1)
}

// ReplaceATag is to replace HTML a tag
func ReplaceATag(html *[]byte, allATags []LinksData, cid string, seed string, uid string) {
	for _, v := range allATags {
		data := url.Values{}
		data.Set("c", cid)
		data.Set("u", uid)
		data.Set("l", v.LinkID)
		data.Set("t", "a")
		hm := campaign.MakeMacSeed(seed, data)

		*html = bytes.Replace(*html, []byte(fmt.Sprintf("href=\"%s\"", v.URL)),
			[]byte(fmt.Sprintf("href=\"https://%s/door/%x?%s\"", os.Getenv("mailbox_web_site"), hm, data.Encode())), -1)
	}
}

// LinksData is the link data
type LinksData struct {
	Md5h   string
	LinkID string
	URL    []byte
}

// FilterATags is to filter, find all a tag data
func FilterATags(body []byte, cid string) []LinksData {
	allATags := areg.FindAllSubmatch(body, -1)
	result := make([]LinksData, len(allATags))
	for i, v := range allATags {
		md5h := md5.New()
		md5h.Write(v[1])
		md5hstr := fmt.Sprintf("%x", md5h.Sum(nil))
		linkID := fmt.Sprintf("%s", utils.GenSeed())
		_, err := utils.GetConn().Query(`INSERT INTO links(id,cid,url,urlhash) VALUES(?,?,?,?)`, linkID, cid, v[1], md5hstr)
		if err != nil {
			rows, _ := utils.GetConn().Query(`SELECT id FROM links WHERE cid=? AND urlhash=?`, cid, md5hstr)
			for rows.Next() {
				rows.Scan(&linkID)
			}
		}
		result[i].Md5h = md5hstr
		result[i].LinkID = linkID
		result[i].URL = v[1]
	}
	return result
}
