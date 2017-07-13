package mails

import (
	"database/sql"
	"fmt"
	"log"
	"os"
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
func GenParams(to string, message string, subject string) *ses.SendEmailInput {
	params := baseParams()
	params.Destination.ToAddresses[0] = aws.String(to)
	params.Message.Body.Html.Data = aws.String(message)
	params.Message.Subject.Data = aws.String(subject)
	return params
}

// Send is to send mail
func Send(params *ses.SendEmailInput) {
	for i := 0; i < 5; i++ {
		resp, err := svc.SendEmail(params)
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
func ProcessSend(body []byte, rows *sql.Rows, cid string, replaceLink bool, subject string, dryRun bool, limit int) {

	var seed = campaign.GetSeed(cid)
	var count int
	var wg sync.WaitGroup
	ql = make(chan struct{}, limit)
	for rows.Next() {
		var (
			allATags map[string]LinksData
			email    string
			fname    string
			lname    string
			msg      []byte
			no       string
		)

		rows.Scan(&no, &email, &fname, &lname)

		msg = body
		ReplaceFname(&msg, fname)
		ReplaceLname(&msg, lname)
		if replaceLink {
			allATags = FilterATags(&msg, cid)
			ReplaceATag(&msg, allATags, cid, seed, no)
		}
		ReplaceReader(&msg, cid, seed, no)
		if dryRun {
			log.Printf("%s\n", msg)
			var n int
			for _, v := range allATags {
				n++
				fmt.Printf("%02d => [%s] %s\n", n, v.LinkID, v.URL)
			}
		} else {
			wg.Add(1)
			ql <- struct{}{}
			go SendWG(GenParams(
				fmt.Sprintf("%s %s <%s>", fname, lname, email),
				string(msg),
				subject), &wg)
		}
		count++
	}
	wg.Wait()
	log.Printf("\n  cid: %s, count: %d\n  Subject: `%s`\n", cid, count, subject)
}
