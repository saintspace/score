package eventpub

import (
	"encoding/json"
	"fmt"
	"net/url"
	"score/app/models"
)

// EventPublisher is responsible for building & publishing events to be processed later by a worker
type EventPublisher struct {
	notifier Notifier
	config   Config
}

type Notifier interface {
	PublishPlatformEventMessage(message string) error
}

type Config interface {
	WebAppDomainName() string
	MainTransactionalSendingAddress() string
}

func New(notifier Notifier, config Config) *EventPublisher {
	return &EventPublisher{
		notifier: notifier,
		config:   config,
	}
}

func (s *EventPublisher) PublishEmailVerificationTask(email, token string) error {
	escapedToken := url.QueryEscape(token)
	linkTemplate := "https://%s/saintspace/universe/verify-email-subscription?token=%s"
	link := fmt.Sprintf(linkTemplate, s.config.WebAppDomainName(), escapedToken)
	emailSendTask := models.EmailSendTaskPlatformEvent{
		TemplateName:  "email-subscription-verification",
		SenderAddress: s.config.MainTransactionalSendingAddress(),
		SubjectLine:   "Confirm Your Subscription",
		ToAddresses:   []string{email},
		Parameters: map[string]string{
			"verificationLink": link,
		},
	}
	emailSendTaskBytes, err := json.Marshal(emailSendTask)
	if err != nil {
		return fmt.Errorf("error while marshaling email send task details => %v", err.Error())
	} else {
		task := models.PlatformEvent{
			EventName:     "email-send-task",
			CorrelationId: "",
			EventDetails:  string(emailSendTaskBytes),
		}
		taskBytes, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("error while marshaling email send task => %v", err.Error())
		} else {
			if err := s.notifier.PublishPlatformEventMessage(string(taskBytes)); err != nil {
				return fmt.Errorf("error while publishing email send task => %v", err.Error())
			}
		}
	}
	return nil
}
