package campaign

import (
	"crypto/hmac"
	"log"
	"net/url"

	"github.com/toomore/mailbox/utils"
)

var cacheSeed map[string]string

func init() {
	cacheSeed = make(map[string]string)
}

// MakeMac is to hmac data from campaign seed
func MakeMac(campaignID string, data url.Values) []byte {
	return utils.GenHmac([]byte(GetSeed(campaignID)), []byte(data.Encode()))
}

// MakeMacSeed is to hmac data with campaign seed
func MakeMacSeed(seed string, data url.Values) []byte {
	return utils.GenHmac([]byte(seed), []byte(data.Encode()))
}

// CheckMac is to check hash mac
func CheckMac(hm []byte, campaignID string, data url.Values) bool {
	return hmac.Equal(hm, MakeMac(campaignID, data))
}

// GetSeed is to get campaign seed
func GetSeed(campaignID string) string {
	var (
		ok   bool
		seed string
	)
	if seed, ok = cacheSeed[campaignID]; ok {
		return seed
	}
	rows, err := utils.GetConn().Query(`SELECT seed FROM campaign WHERE id=? `, campaignID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		rows.Scan(&seed)
	}
	if seed == "" {
		log.Fatalln("Find no campaign ID")
	}
	cacheSeed[campaignID] = seed
	return seed
}
