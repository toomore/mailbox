package reader

import "github.com/toomore/mailbox/utils"

// Save is to save read record
func Save(cid, uid, ip, agent string) {
	conn := utils.GetConn()
	conn.Query(`INSERT INTO reader(cid,uid,ip,agent) VALUES(?,?,?,?)`, cid, uid, ip, agent)
}
