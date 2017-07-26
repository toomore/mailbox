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
	html := []byte(`<a href="https://toomore.net/">1</a><a href="{{WASHI}}https://toomore.net/?name={{FNAME}}{{/WASHI}}">`)
	links := FilterATags(&html, "12345678")
	t.Logf("%+v\n", links)
	if links["585d7c77"].Md5h != "94fab8b14a0214468e705a9cbcfd68c2" {
		t.Error("Should be `94fab8b14a0214468e705a9cbcfd68c2`")
	}

	ReplaceATag(&html, links, "12345678", "87654321", "11")
	t.Logf("%s", html)

	links = FilterWashiTags(&html, "12345678")
	t.Logf("%+v\n", links)
	if links["83fcc9b2"].Md5h != "239b380202a2c8c65df27aed9135fe5c" {
		t.Error("Should be `239b380202a2c8c65df27aed9135fe5c`")
	}
	ReplaceWashiTag(&html, links, "12345678", "87654321", "11")
	t.Logf("%s", html)
	if fmt.Sprintf("%s", html) != `<a href="https://example.com/door/673741b4a6af077e0a9f23fa1e1b0fe0c6c74c68682a6e9f983d281c5d70f81c?c=12345678&l=585d7c77&t=a&u=11">1</a><a href="https://example.com/washi/645c9c7185b373b38f68011f246ddfe5c109ed37bb7c671aaac0dc9b1b4c89ab?c=12345678&l=83fcc9b2&u=11">` {
		t.Error("Replace fail")
	}
}
