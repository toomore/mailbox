package mails

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	os.Setenv("mailbox_web_site", "example.com")
}

func TestReplaceReader(t *testing.T) {
	html := []byte(`<img src="{{READER}}">`)
	ReplaceReader(&html, "12345678", "87654321", "11")
	if fmt.Sprintf("%s", html) != `<img src="https://example.com/read/f1321a88e777a754e1a4b6121a495331f5c57c6cd477cec6db3bed4cff71c6cb?c=12345678&u=11">` {
		t.Error(`Should be <img src="https://example.com/read/f1321a88e777a754e1a4b6121a495331f5c57c6cd477cec6db3bed4cff71c6cb?c=12345678&u=11">`)
	}
}

func TestReplaceFname(t *testing.T) {
	html := []byte(`{{FNAME}} {{LNAME}}`)
	ReplaceFname(&html, "Toomore")
	if fmt.Sprintf("%s", html) != "Toomore {{LNAME}}" {
		t.Error("Should be `Toomore {{LANME}}`")
	}
	ReplaceLname(&html, "Chiang")
	if fmt.Sprintf("%s", html) != "Toomore Chiang" {
		t.Error("Should be `Toomore Chiang`")
	}
}

func TestFilterATags(t *testing.T) {
	html := []byte(`<a href="https://toomore.net/">1</a><a href="{{WASHI}}https://toomore.net/?name={{FNAME}}{{/WASHI}}">2</a>`)
	links := FilterATags(&html, "12345678")
	t.Logf("%+v\n", links)

	for _, v := range links {
		if v.Md5h != "94fab8b14a0214468e705a9cbcfd68c2" {
			t.Error("Should be `94fab8b14a0214468e705a9cbcfd68c2`")
		}
	}

	ReplaceATag(&html, links, "12345678", "87654321", "11")
	t.Logf("%s", html)

	links = FilterWashiTags(&html, "12345678")
	t.Logf("%+v\n", links)

	for _, v := range links {
		if v.Md5h != "239b380202a2c8c65df27aed9135fe5c" {
			t.Error("Should be `239b380202a2c8c65df27aed9135fe5c`")
		}
	}
	ReplaceWashiTag(&html, links, "12345678", "87654321", "11")
	t.Logf("%s", html)
}
