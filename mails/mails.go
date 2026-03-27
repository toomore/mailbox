package mails

import (
	"database/sql"
	"fmt"
	"log"
	"mime"
	"net/mail"
	"os"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/toomore/mailbox/campaign"
)

var svc *ses.SES
var ql chan struct{}

func init() {
	svc = ses.New(session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("mailbox_ses_key"), os.Getenv("mailbox_ses_token"), ""),
		},
	),
	))
}

func baseParams() *ses.SendEmailInput {
	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String("")},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data:    aws.String(""),
					Charset: aws.String("UTF-8"),
				},
				Text: &ses.Content{
					Data:    aws.String(""),
					Charset: aws.String("UTF-8"),
				},
			},
			Subject: &ses.Content{
				Data:    aws.String(""),
				Charset: aws.String("UTF-8"),
			},
		},
		Source: aws.String(os.Getenv("mailbox_ses_sender")),
	}
}

// GenParams is to gen email params
func GenParams(to string, message string, subject string, text string) *ses.SendEmailInput {
	params := baseParams()
	params.Destination.ToAddresses[0] = aws.String(to)
	params.Message.Body.Html.Data = aws.String(message)
	params.Message.Body.Text.Data = aws.String(text)
	params.Message.Subject.Data = aws.String(subject)
	if os.Getenv("mailbox_ses_replyto") != "" {
		params.ReplyToAddresses = []*string{aws.String(os.Getenv("mailbox_ses_replyto"))}
	}
	return params
}

func getUnsubscribeMailto() string {
	if v := strings.TrimSpace(os.Getenv("mailbox_unsubscribe_mailto")); v != "" {
		return v
	}
	return strings.TrimSpace(os.Getenv("mailbox_ses_replyto"))
}

func buildListUnsubscribeHeaders() []string {
	mailto := getUnsubscribeMailto()
	if mailto == "" {
		return nil
	}
	headers := []string{
		fmt.Sprintf("List-Unsubscribe: <mailto:%s?subject=unsubscribe>", mailto),
	}
	if strings.EqualFold(strings.TrimSpace(os.Getenv("mailbox_unsubscribe_one_click")), "true") ||
		strings.TrimSpace(os.Getenv("mailbox_unsubscribe_one_click")) == "1" {
		headers = append(headers, "List-Unsubscribe-Post: List-Unsubscribe=One-Click")
	}
	return headers
}

func encodeSubject(subject string) string {
	if subject == "" {
		return ""
	}
	return mime.BEncoding.Encode("UTF-8", subject)
}

func buildRawEmail(params *ses.SendEmailInput) []byte {
	var (
		to      = strings.TrimSpace(aws.StringValue(params.Destination.ToAddresses[0]))
		from    = strings.TrimSpace(aws.StringValue(params.Source))
		subject = aws.StringValue(params.Message.Subject.Data)
		html    = aws.StringValue(params.Message.Body.Html.Data)
		text    = aws.StringValue(params.Message.Body.Text.Data)
	)
	boundary := "mailbox_alt_boundary"
	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", to),
		fmt.Sprintf("Subject: %s", encodeSubject(subject)),
		"MIME-Version: 1.0",
		fmt.Sprintf("Content-Type: multipart/alternative; boundary=%q", boundary),
	}
	if len(params.ReplyToAddresses) > 0 {
		reply := strings.TrimSpace(aws.StringValue(params.ReplyToAddresses[0]))
		if reply != "" {
			headers = append(headers, fmt.Sprintf("Reply-To: %s", reply))
		}
	}
	headers = append(headers, buildListUnsubscribeHeaders()...)

	var msg strings.Builder
	msg.WriteString(strings.Join(headers, "\r\n"))
	msg.WriteString("\r\n\r\n")
	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
	msg.WriteString(text)
	msg.WriteString("\r\n")
	msg.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("Content-Transfer-Encoding: 8bit\r\n\r\n")
	msg.WriteString(html)
	msg.WriteString("\r\n")
	msg.WriteString(fmt.Sprintf("--%s--\r\n", boundary))
	return []byte(msg.String())
}

// Send is to send mail
func Send(params *ses.SendEmailInput) {
	// Validate common address formats before passing to SES raw API.
	to := strings.TrimSpace(aws.StringValue(params.Destination.ToAddresses[0]))
	from := strings.TrimSpace(aws.StringValue(params.Source))
	if _, err := mail.ParseAddress(to); err != nil {
		log.Println(to, err)
		return
	}
	if _, err := mail.ParseAddress(from); err != nil {
		log.Println(from, err)
		return
	}
	raw := buildRawEmail(params)
	for i := 0; i < 5; i++ {
		resp, err := svc.SendRawEmail(&ses.SendRawEmailInput{
			RawMessage: &ses.RawMessage{Data: raw},
		})
		if err == nil {
			log.Println(*params.Destination.ToAddresses[0], resp)
			return
		}
		log.Println(*params.Destination.ToAddresses[0], err)
	}
}

// SendWG is with WaitGroup
func SendWG(params *ses.SendEmailInput, wg *sync.WaitGroup) {
	Send(params)
	<-ql
	wg.Done()
}

// ProcessSend is to start send from rows
func ProcessSend(body []byte, text []byte, rows *sql.Rows, cid string, replaceLink bool, subject string, dryRun bool, limit int) {

	var seed = campaign.GetSeed(cid)
	var count int
	var wg sync.WaitGroup
	ql = make(chan struct{}, limit)
	for rows.Next() {
		var (
			allATags     map[string]LinksData
			allWashiTags map[string]LinksData
			email        string
			fname        string
			lname        string
			msg          []byte
			msg_text     []byte
			no           string
			subjectbyte  []byte
		)

		rows.Scan(&no, &email, &fname, &lname)

		msg = body
		msg_text = text
		if replaceLink {
			allATags = FilterATags(&msg, cid)
			ReplaceATag(&msg, allATags, cid, seed, no)
			allWashiTags = FilterWashiTags(&msg, cid)
			ReplaceWashiTag(&msg, allWashiTags, cid, seed, no)

			allATags = FilterATags(&msg_text, cid)
			ReplaceATag(&msg_text, allATags, cid, seed, no)
			allWashiTags = FilterWashiTags(&msg_text, cid)
			ReplaceWashiTag(&msg_text, allWashiTags, cid, seed, no)
		}
		ReplaceFname(&msg, fname)
		ReplaceLname(&msg, lname)
		ReplaceReader(&msg, cid, seed, no)

		ReplaceFname(&msg_text, fname)
		ReplaceLname(&msg_text, lname)
		ReplaceReader(&msg_text, cid, seed, no)

		subjectbyte = []byte(subject)
		ReplaceFname(&subjectbyte, fname)
		ReplaceLname(&subjectbyte, lname)

		if dryRun {
			log.Printf("%s\n", msg)
			log.Printf("%s\n", msg_text)
			var n int
			for _, v := range allATags {
				n++
				fmt.Printf("%02d => [%s] %s\n", n, v.LinkID, v.URL)
			}
			for _, v := range allWashiTags {
				n++
				fmt.Printf("%02d => [W][%s] %s\n", n, v.LinkID, v.URL)
			}
			fmt.Printf("Subject: %s\n", subjectbyte)
		} else {
			wg.Add(1)
			ql <- struct{}{}
			go SendWG(GenParams(
				fmt.Sprintf("%s %s <%s>", fname, lname, email),
				string(msg),
				string(subjectbyte),
				string(msg_text),
			), &wg)
		}
		count++
	}
	wg.Wait()
	log.Printf("\n  cid: %s, count: %d\n  Subject: `%s`\n", cid, count, subject)
}
