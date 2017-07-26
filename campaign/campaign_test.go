package campaign

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestMakeMac(t *testing.T) {
	campaignID, _ := Create()
	cid := fmt.Sprintf("%s", campaignID)

	data := url.Values{}
	data.Set("name", "toomore")
	data.Set("age", "30")
	data.Set("cid", cid)
	t.Log(data.Encode())
	t.Logf("%x\n", MakeMac(cid, data))
}

func TestCheckMac(t *testing.T) {
	campaignID, seed := Create()
	cid := fmt.Sprintf("%s", campaignID)
	t.Logf("cid: %s, seed: %s", campaignID, seed)

	data := url.Values{}
	data.Set("name", "toomore")
	data.Set("age", "30")
	data.Set("cid", cid)
	t.Log(data.Encode())

	hm := MakeMac(cid, data)
	t.Logf("%x\n", hm)
	t.Log(CheckMac(hm, cid, data))
}

func ExampleCheckMac() {
	hm := "fd3890c84e29acf04e1bc1cd2c6d37d49931209669bf94cacad65540570d9c12"
	data := url.Values{}

	// If hmac value from string, need using `hex.DecodeString` to byte
	hmbyte, _ := hex.DecodeString(hm)
	CheckMac(hmbyte, "123", data)
}
