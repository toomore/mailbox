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

func TestMakeMacSeed(t *testing.T) {
	data := url.Values{}
	data.Set("name", "toomore")
	data.Set("age", "30")
	data.Set("cid", "12345678")
	if fmt.Sprintf("%x", MakeMacSeed("87654321", data)) != "37f90a0e96cebc981fa72cc64a79a942374906367dbc8c0fce44f9ec6dd7fe27" {
		t.Error("Should be `37f90a0e96cebc981fa72cc64a79a942374906367dbc8c0fce44f9ec6dd7fe27`")
	}
}

func ExampleCheckMac() {
	hm := "fd3890c84e29acf04e1bc1cd2c6d37d49931209669bf94cacad65540570d9c12"
	data := url.Values{}

	// If hmac value from string, need using `hex.DecodeString` to byte
	hmbyte, _ := hex.DecodeString(hm)
	CheckMac(hmbyte, "123", data)
}
