package campaign

import (
	"crypto/hmac"
	"log"
	"net/url"

	"github.com/toomore/mailbox/utils"
)

// MakeMac is to hmac data with campaign seed
func MakeMac(campaignID string, data url.Values) []byte {
	conn := utils.GetConn()
	rows, err := conn.Query(`SELECT seed FROM campaign WHERE id='?' `, campaignID)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var seed string
	for rows.Next() {
		rows.Scan(&seed)
	}
	return utils.GenHmac([]byte(seed), []byte(data.Encode()))
}

// CheckMac is to check hash mac
func CheckMac(hm []byte, campaignID string, data url.Values) bool {
	return hmac.Equal(hm, MakeMac(campaignID, data))
}
