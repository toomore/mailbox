package mails

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var svc *ses.SES

func init() {
	svc = ses.New(session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("mailbox_ses_api"), os.Getenv("mailbox_ses_key"), ""),
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
