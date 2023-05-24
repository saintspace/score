package sns

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/google/uuid"
)

// SNS is responsible for interfacing with AWS Simple Notification Service
type SNS struct {
	svc    *sns.SNS
	config Config
}

func New(awsSession *session.Session, config Config) *SNS {
	snsSession := sns.New(awsSession)
	return &SNS{
		svc:    snsSession,
		config: config,
	}
}

type Config interface {
	PlatformEventsTopicArn() string
}

func (s *SNS) publishMessageToTopic(message string, topicArn string) error {
	uniqueMessageId := uuid.New().String()
	_, err := s.svc.Publish(&sns.PublishInput{
		Message:                &message,
		TopicArn:               &topicArn,
		MessageGroupId:         &uniqueMessageId,
		MessageDeduplicationId: &uniqueMessageId,
	})
	return err
}

func (s *SNS) PublishPlatformEventMessage(serializedPlatformEvent string) error {
	return s.publishMessageToTopic(serializedPlatformEvent, s.config.PlatformEventsTopicArn())
}
