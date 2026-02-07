package reader

import (
	"log"

	"github.com/toomore/mailbox/utils"
)

// Save is to save read record
func Save(cid, uid, ip, agent string) {
	_, err := utils.GetConn().Exec(`INSERT INTO reader(cid,uid,ip,agent) VALUES(?,?,?,?)`, cid, uid, ip, agent)
	if err != nil {
		log.Println("[reader][Save] ", err)
	}
}
