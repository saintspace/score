package ses

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// SES is responsible for interfacing with AWS Simple Email Service
type SES struct {
	svc *ses.SES
}

func New(awsSession *session.Session) *SES {
	sesSession := ses.New(awsSession)
	return &SES{
		svc: sesSession,
	}
}

func (s *SES) SendEmail(
	senderAddress string,
	toAddresses []string,
	subjectLine string,
	emailBody string,
) error {
	email := &ses.SendEmailInput{
		Source: aws.String(senderAddress),
		Destination: &ses.Destination{
			ToAddresses: aws.StringSlice(toAddresses),
		},
		Message: &ses.Message{
			Subject: &ses.Content{Data: aws.String(subjectLine)},
			Body:    &ses.Body{Text: &ses.Content{Data: aws.String(emailBody)}},
		},
	}
	_, err := s.svc.SendEmail(email)
	return err
}
