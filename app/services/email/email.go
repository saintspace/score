package email

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/mail"
	"score/app/models"
	"strings"
)

type EmailService struct {
	emailSender    EmailSender
	datastore      Datastore
	logger         Logger
	eventPublisher PlatformEventPublisher
}

func New(emailSender EmailSender, datastore Datastore, logger Logger, eventPublisher PlatformEventPublisher) *EmailService {
	return &EmailService{
		emailSender:    emailSender,
		datastore:      datastore,
		logger:         logger,
		eventPublisher: eventPublisher,
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

type PlatformEventPublisher interface {
	PublishEmailVerificationTask(email, token string) error
}

type Datastore interface {
	EmailSubscriptionExists(email string) (bool, error)
	AddComplaintToEmailSubscription(email, complaintDetails string, complaintDateUnix int64) error
	AddBounceToEmailSubscription(email, bounceType, bounceDetails string, bounceDateUnix int64) error
	GetEmailSubscription(email string) (*models.EmailSubscription, error)
	CreateEmailSubscription(email string, subscriptionToken string, isVerified bool) error
	VerifyEmailSubscription(email string) error
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
		exists, err := s.datastore.EmailSubscriptionExists(email)
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
		exists, err := s.datastore.EmailSubscriptionExists(email)
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

func (s *EmailService) IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (s *EmailService) EmailSubscriptionExists(email string) (bool, error) {
	return s.datastore.EmailSubscriptionExists(email)
}

func (s *EmailService) CreateEmailSubscription(email string) error {
	subscriptionToken, err := generateEmailSubscriptionToken(email)
	if err != nil {
		return fmt.Errorf("error generating email subscription token: %v", err)
	}
	err = s.datastore.CreateEmailSubscription(email, subscriptionToken, false)
	if err != nil {
		return fmt.Errorf("error creating email subscription in datastore: %v", err)
	}
	err = s.eventPublisher.PublishEmailVerificationTask(email, subscriptionToken)
	if err != nil {
		return fmt.Errorf("error publishing email verification task: %v", err)
	}
	return nil
}

func (s *EmailService) VerifyEmailWithSubscriptionToken(token string) error {
	email, err := parseEmailFromSubscriptionToken(token)
	if err != nil {
		return fmt.Errorf("error parsing email from token: %v", err)
	}
	subscriptionExists, err := s.datastore.EmailSubscriptionExists(email)
	if err != nil {
		return fmt.Errorf("error checking if email subscription exists: %v", err)
	}
	if !subscriptionExists {
		return fmt.Errorf("email subscription doesn't exist: %v", email)
	}
	err = s.datastore.VerifyEmailSubscription(email)
	if err != nil {
		return fmt.Errorf("error verifying existing subscription: %v", err)
	}
	return nil
}

func generateEmailSubscriptionToken(email string) (string, error) {
	// Generate a 32-byte random token
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	// Encode the token and email address as a base64 string
	token := base64.StdEncoding.EncodeToString(tokenBytes)
	emailToken := base64.StdEncoding.EncodeToString([]byte(email))

	// Add the email token to the token as a prefix, separated by a colon
	tokenWithPrefix := emailToken + ":" + token

	return tokenWithPrefix, nil
}

func parseEmailFromSubscriptionToken(tokenWithPrefix string) (string, error) {
	// Split the token into email token and token parts
	parts := strings.Split(tokenWithPrefix, ":")
	if len(parts) != 2 {
		return "", errors.New("invalid token")
	}

	// Decode the email token and return the email address
	emailBytes, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}

	return string(emailBytes), nil
}
