package mails

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sync"

	"github.com/toomore/mailbox/campaign"
	"github.com/toomore/mailbox/utils"
)

var (
	htmla    = regexp.MustCompile(`href="(http[s]?://[a-zA-z0-9/\.:?=,-@%()_&\+#]+)"`)
	washireg = regexp.MustCompile(`href="({{WASHI}}.+{{\/WASHI}})"`)
)

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

// ReplaceLname is to replace FNAME tag
func ReplaceLname(html *[]byte, lname string) {
	*html = bytes.Replace(*html, []byte("{{LNAME}}"), []byte(lname), -1)
}

// ReplaceATag is to replace HTML a tag
func ReplaceATag(html *[]byte, allATags map[string]LinksData, cid string, seed string, uid string) {
	data := url.Values{}
	data.Set("c", cid)
	data.Set("u", uid)
	data.Set("t", "a")

	for _, v := range allATags {
		data.Set("l", v.LinkID)
		hm := campaign.MakeMacSeed(seed, data)

		*html = bytes.Replace(*html, []byte(fmt.Sprintf("href=\"%s\"", v.URL)),
			[]byte(fmt.Sprintf("href=\"https://%s/door/%x?%s\"", os.Getenv("mailbox_web_site"), hm, data.Encode())), -1)
	}
}

// ReplaceWashiTag is to replace HTML a tag
func ReplaceWashiTag(html *[]byte, allATags map[string]LinksData, cid string, seed string, uid string) {
	data := url.Values{}
	data.Set("c", cid)
	data.Set("u", uid)

	for _, v := range allATags {
		data.Set("l", v.LinkID)
		hm := campaign.MakeMacSeed(seed, data)

		*html = bytes.Replace(*html, []byte(fmt.Sprintf("href=\"%s\"", v.URL)),
			[]byte(fmt.Sprintf("href=\"https://%s/washi/%x?%s\"", os.Getenv("mailbox_web_site"), hm, data.Encode())), -1)
	}
}

// LinksData is the link data
type LinksData struct {
	Md5h   string
	LinkID string
	URL    []byte
}

// FilterATags is to filter, find all a tag data
func FilterATags(body *[]byte, cid string) map[string]LinksData {
	return filteratags(htmla, body, cid)
}

// FilterWashiTags is to filter, find all {{WASHI}} tag data
func FilterWashiTags(body *[]byte, cid string) map[string]LinksData {
	return filteratags(washireg, body, cid)
}

func filteratags(rg *regexp.Regexp, body *[]byte, cid string) map[string]LinksData {
	var (
		allATags = rg.FindAllSubmatch(*body, -1)
		conn     = utils.GetConn()
		result   = make(map[string]LinksData)
		wg       sync.WaitGroup
		lock     = sync.RWMutex{}
	)
	wg.Add(len(allATags))
	for _, v := range allATags {
		go func(url []byte) {
			md5h := md5.New()
			md5h.Write(url)
			md5hstr := fmt.Sprintf("%x", md5h.Sum(nil))
			linkID := fmt.Sprintf("%x", utils.GenSeed())
			_, err := conn.Query(`INSERT INTO links(id,cid,url,urlhash) VALUES(?,?,?,?)`, linkID, cid, url, md5hstr)
			if err != nil {
				rows, _ := conn.Query(`SELECT id FROM links WHERE cid=? AND urlhash=?`, cid, md5hstr)
				for rows.Next() {
					rows.Scan(&linkID)
				}
			}
			lock.Lock()
			result[linkID] = LinksData{
				Md5h:   md5hstr,
				LinkID: linkID,
				URL:    url,
			}
			lock.Unlock()
			wg.Done()
		}(v[1])
	}
	wg.Wait()
	conn.Close()
	return result
}
