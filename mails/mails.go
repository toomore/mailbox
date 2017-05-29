package mails

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/toomore/mailbox/campaign"
)

var svc *ses.SES

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
	resp, err := svc.SendEmail(params)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(resp)
	}
}

// ProcessSend is to start send from rows
func ProcessSend(body []byte, rows *sql.Rows, cid string, replaceLink bool, subject string, dryRun bool) {
	var allATags []LinksData
	if replaceLink {
		allATags = FilterATags(body, cid)
	}

	var seed = campaign.GetSeed(cid)
	var count int
	for rows.Next() {
		var (
			email string
			fname string
			lname string
			msg   []byte
			no    string
		)
		rows.Scan(&no, &email, &fname, &lname)

		msg = body
		if replaceLink {
			ReplaceATag(&msg, allATags, cid, seed, no)
		}
		ReplaceFname(&msg, fname)
		ReplaceReader(&msg, cid, seed, no)
		if dryRun {
			log.Printf("%s\n", msg)
			for i, v := range allATags {
				fmt.Printf("%02d => [%s] %s\n", i, v.LinkID, v.URL)
			}
		} else {
			Send(GenParams(
				fmt.Sprintf("%s %s <%s>", fname, lname, email),
				string(msg),
				subject))
		}
		count++
	}
	log.Printf("\n  cid: %s, count: %d\n  Subject: `%s`\n", cid, count, subject)
}
