package email

import (
	"fmt"
	"score/app/models"
)

type EmailService struct {
	emailSender EmailSender
	datastore   Datastore
	logger      Logger
}

func New(emailSender EmailSender, datastore Datastore, logger Logger) *EmailService {
	return &EmailService{
		emailSender: emailSender,
		datastore:   datastore,
		logger:      logger,
	}
}

type EmailSender interface {
	SendEmail(
		senderAddress string,
		toAddresses []string,
		subjectLine string,
		emailBody string,
	) error
}

type Datastore interface {
	EmailSubscriptionItemExists(email string) (bool, error)
	AddComplaintToEmailSubscription(email, complaintDetails string, complaintDateUnix int64) error
	AddBounceToEmailSubscription(email, bounceType, bounceDetails string, bounceDateUnix int64) error
	GetEmailSubscription(email string) (*models.EmailSubscription, error)
}

type Logger interface {
	InfoWithContext(message string, keysAndValues ...interface{})
	ErrorWithContext(message string, keysAndValues ...interface{})
	DebugWithContext(message string, keysAndValues ...interface{})
}

func (s *EmailService) SendTemplatedEmail(
	templateName string,
	templateParams map[string]string,
	subjectLine string,
	senderAddress string,
	toAddresses []string,
) error {
	// Get email subscriptions and filter out email addresses that have complaints or bounces
	filteredToAddresses := []string{}
	for _, email := range toAddresses {
		emailSubscription, err := s.datastore.GetEmailSubscription(email)
		if err != nil {
			s.logger.ErrorWithContext("error getting email subscription", "email", email, "error", err.Error())
			return err
		} else if emailSubscription == nil {
			s.logger.InfoWithContext("email subscription not found. excluding from toAddresses", "email", email)
			continue
		} else if emailSubscription.HasComplaint {
			s.logger.InfoWithContext("email subscription has complaint. excluding from toAddresses", "email", email)
			continue
		} else if emailSubscription.BounceType == "Permanent" {
			s.logger.InfoWithContext("email subscription has permanent bounce. excluding from toAddresses", "email", email)
			continue
		} else {
			filteredToAddresses = append(filteredToAddresses, email)
		}
	}
	if len(filteredToAddresses) == 0 {
		s.logger.InfoWithContext("no email addresses to send to after filtering", "templateName", templateName)
		return nil
	}
	emailBody, err := generateEmailBodyFromTemplate(templateName, templateParams)
	if err != nil {
		return fmt.Errorf("error while generating email body => %v", err.Error())
	}
	return s.emailSender.SendEmail(senderAddress, filteredToAddresses, subjectLine, emailBody)
}

func (s *EmailService) ProcessEmailComplaint(
	complainedEmailAddresses []string,
	complaintDetails string,
	complaintUnixTime int64,
) error {
	for _, email := range complainedEmailAddresses {
		exists, err := s.datastore.EmailSubscriptionItemExists(email)
		if err != nil {
			return fmt.Errorf("error checking if email subscription item exists => %v", err.Error())
		}
		if exists {
			err := s.datastore.AddComplaintToEmailSubscription(
				email,
				complaintDetails,
				complaintUnixTime,
			)
			if err != nil {
				return fmt.Errorf("error while adding complaint to email subscription => %v", err.Error())
			}
		}
	}
	return nil
}

// ProcessEmailBounce processes an email bounce event
func (s *EmailService) ProcessEmailBounce(
	bouncedEmailAddresses []string,
	bounceType,
	bounceSubType,
	bounceDetails string,
	bounceUnixTime int64,
) error {
	for _, email := range bouncedEmailAddresses {
		exists, err := s.datastore.EmailSubscriptionItemExists(email)
		if err != nil {
			return fmt.Errorf("error checking if email subscription item exists => %v", err.Error())
		}
		if exists {
			savedBounceType := ""
			if bounceType == "Permanent" {
				savedBounceType = "Permanent"
			} else {
				savedBounceType = bounceSubType
			}
			err := s.datastore.AddBounceToEmailSubscription(
				email,
				savedBounceType,
				bounceDetails,
				bounceUnixTime,
			)
			if err != nil {
				return fmt.Errorf("error while adding bounce to email subscription => %v", err.Error())
			}
		}
	}
	return nil
}
