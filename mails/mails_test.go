package mails

import (
	"fmt"
	"testing"

	"github.com/toomore/mailbox/campaign"
	"github.com/toomore/mailbox/utils"
)

func TestProcessSend(t *testing.T) {
	stmt, err := utils.GetConn().Prepare(`INSERT INTO user(email,groups,f_name,l_name)
	                           VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE f_name=?, l_name=?`)
	if err != nil {
		t.Fatal(err)
	}
	stmt.Exec("toomore0929@gmail.com", "test", "Toomore", "Chiang", "Toomore", "Chiang")
	rows, err := utils.GetConn().Query("select id,email,f_name,l_name from user")

	if err != nil {
		t.Fatal(err)
	}
	cid, _ := campaign.Create()
	ProcessSend(
		[]byte(`<a href="https://toomore.net/">1</a><a href="{{WASHI}}https://toomore.net/{{/WASHI}}">2</a>`),
		rows,
		fmt.Sprintf("%s", cid),
		true,
		"Test",
		true,
		4)
}
