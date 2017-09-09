package campaign

import (
	"crypto/hmac"
	"fmt"
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
		log.Fatal("[campaign][GetSeed] ", err)
	}
	for rows.Next() {
		rows.Scan(&seed)
	}
	if seed == "" {
		log.Fatal("[campaign][GetSeed] Find no campaign ID")
	}
	cacheSeed[campaignID] = seed
	return seed
}

// Create is to create a new campaign id and seed
func Create() ([]byte, []byte) {
	id, seed := utils.GenSeed(), utils.GenSeed()
	rows, err := utils.GetConn().Query(fmt.Sprintf(`INSERT INTO campaign(id,seed) VALUES('%s', '%s')`, id, seed))
	defer rows.Close()

	if err != nil {
		log.Fatal("[campaign][Create] ", err)
	}

	return id, seed
}
