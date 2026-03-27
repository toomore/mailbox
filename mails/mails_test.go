package mails

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/toomore/mailbox/campaign"
	"github.com/toomore/mailbox/utils"
)

func TestGenParams(t *testing.T) {
	t.Logf("%+v", GenParams("toomore0929@gmail.com", "message", "[Test]", "text message"))
}

func TestBuildListUnsubscribeHeaders(t *testing.T) {
	t.Setenv("mailbox_unsubscribe_mailto", "sender+unsubscribe@example.com")
	t.Setenv("mailbox_unsubscribe_one_click", "true")
	headers := buildListUnsubscribeHeaders()
	if len(headers) != 2 {
		t.Fatalf("unexpected header count: %d", len(headers))
	}
	if headers[0] != "List-Unsubscribe: <mailto:sender+unsubscribe@example.com?subject=unsubscribe>" {
		t.Fatalf("unexpected List-Unsubscribe header: %s", headers[0])
	}
	if headers[1] != "List-Unsubscribe-Post: List-Unsubscribe=One-Click" {
		t.Fatalf("unexpected List-Unsubscribe-Post header: %s", headers[1])
	}
}

func TestBuildRawEmailContainsHeaders(t *testing.T) {
	t.Setenv("mailbox_ses_sender", "Sender <sender@example.com>")
	t.Setenv("mailbox_ses_replyto", "sender+reply@example.com")
	t.Setenv("mailbox_unsubscribe_mailto", "sender+unsubscribe@example.com")
	t.Setenv("mailbox_unsubscribe_one_click", "1")
	params := GenParams("User <user@example.com>", "<b>Hello</b>", "測試 Subject", "Hello")
	raw := string(buildRawEmail(params))
	for _, expected := range []string{
		"Reply-To: sender+reply@example.com",
		"List-Unsubscribe: <mailto:sender+unsubscribe@example.com?subject=unsubscribe>",
		"List-Unsubscribe-Post: List-Unsubscribe=One-Click",
		"Content-Type: multipart/alternative",
	} {
		if !strings.Contains(raw, expected) {
			t.Fatalf("raw email missing %q", expected)
		}
	}
}

func TestGetUnsubscribeMailtoFallbackReplyTo(t *testing.T) {
	os.Unsetenv("mailbox_unsubscribe_mailto")
	t.Setenv("mailbox_ses_replyto", "sender+unsubscribe@example.com")
	got := getUnsubscribeMailto()
	if got != "sender+unsubscribe@example.com" {
		t.Fatalf("unexpected fallback value: %s", got)
	}
}

func TestProcessSend(t *testing.T) {
	stmt, err := utils.GetConn().Prepare(`INSERT INTO user(email,email_uni,groups,f_name,l_name)
	                           VALUES(?,?,?,?,?) ON DUPLICATE KEY UPDATE f_name=?, l_name=?`)
	if err != nil {
		t.Fatal(err)
	}
	stmt.Exec("toomore0929+123@gmail.com", "toomore0929@gmail.com", "test", "Toomore", "Chiang", "Toomore", "Chiang")
	rows, err := utils.GetConn().Query("select id,email,f_name,l_name from user where groups='test'")

	if err != nil {
		t.Fatal(err)
	}
	cid, _ := campaign.Create()
	// Test with dry run
	ProcessSend(
		[]byte(`<a href="https://toomore.net/">1</a><a href="{{WASHI}}https://toomore.net/{{/WASHI}}">2</a>`),
		[]byte(`<https://toomore.net> {{WASHI}}https://toomore.net/{{/WASHI}}`),
		rows,
		fmt.Sprintf("%x", cid),
		true,
		"Test",
		true,
		4)

	stmt.Exec("to", "test2", "Toomore", "Chiang", "Toomore", "Chiang")
	rows, err = utils.GetConn().Query("select id,email,email_uni,f_name,l_name from user where groups='test2'")

	if err != nil {
		t.Fatal(err)
	}
	// Test Run
	ProcessSend(
		[]byte(`<a href="https://toomore.net/">1</a><a href="{{WASHI}}https://toomore.net/{{/WASHI}}">2</a>`),
		[]byte(`<https://toomore.net> {{WASHI}}https://toomore.net/{{/WASHI}}`),
		rows,
		fmt.Sprintf("%x", cid),
		true,
		"Test",
		false,
		4)
}
