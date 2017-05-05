package campaign

import (
	"net/url"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMakeMac(t *testing.T) {
	data := url.Values{}
	data.Set("name", "toomore")
	data.Set("age", "30")
	data.Set("cid", "f2d024db")
	t.Log(data.Encode())
	t.Logf("%x\n", MakeMac("f2d024db", data))
}

func TestCheckMac(t *testing.T) {
	campaignID := "f2d024db"
	data := url.Values{}
	data.Set("name", "toomore")
	data.Set("age", "30")
	data.Set("cid", campaignID)
	t.Log(data.Encode())
	hm := MakeMac(campaignID, data)
	t.Logf("%x\n", hm)
	t.Log(CheckMac(hm, campaignID, data))
}
