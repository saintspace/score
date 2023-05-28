package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"score/app/models"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Worker struct {
	emailService EmailService
	userService  UserService
	logger       Logger
}

func New(emailService EmailService, userService UserService, logger Logger) *Worker {
	return &Worker{
		emailService: emailService,
		userService:  userService,
		logger:       logger,
	}
}

type EmailService interface {
	SendTemplatedEmail(
		templateName string,
		templateParams map[string]string,
		subjectLine string,
		senderAddress string,
		toAddresses []string,
	) error
	ProcessEmailComplaint(complainedEmailAddresses []string, complaintDetails string, complaintUnixTime int64) error
	ProcessEmailBounce(bouncedEmailAddresses []string, bounceType, bounceSubType, bounceDetails string, bounceUnixTime int64) error
}

type UserService interface {
	CreateUser(email, cognitoUserName string) error
}

type Logger interface {
	InfoWithContext(message string, keysAndValues ...interface{})
	ErrorWithContext(message string, keysAndValues ...interface{})
	DebugWithContext(message string, keysAndValues ...interface{})
}

func (s *Worker) ProcessEvent(eventString string) error {
	event := models.PlatformEvent{}
	err := json.Unmarshal([]byte(eventString), &event)
	if err != nil {
		return fmt.Errorf("error parsing event: %s", err.Error())
	}
	if event.EventName == models.PlatformEventNames.EmailBounce || event.EventName == models.PlatformEventNames.EmailComplaint {
		event.EventDetails = eventString
	}
	switch event.EventName {
	case models.PlatformEventNames.EmailSendTask:
		err = s.HandleEmailSendTask(event.EventDetails)
	case models.PlatformEventNames.AccountConfirmationTask:
		err = s.HandleAccountConfirmationTask(event.EventDetails)
	case models.PlatformEventNames.EmailBounce:
		err = s.HandleEmailBounce(event.EventDetails)
	case models.PlatformEventNames.EmailComplaint:
		err = s.HandleEmailComplaint(event.EventDetails)
	default:
		err = fmt.Errorf("unsupported event name")
	}
	if err != nil {
		s.logger.ErrorWithContext("error processing event", "eventName", event.EventName, "errorMessage", err.Error(), "correlationId", event.CorrelationId)
		return err
	}
	s.logger.InfoWithContext("successfully processed event", "eventName", event.EventName, "eventDetails", event.EventDetails, "correlationId", event.CorrelationId)
	return nil
}

func (s *Worker) Run() error {
	var lambdaEventHandler = func(ctx context.Context, sqsEvent events.SQSEvent) error {
		var wg sync.WaitGroup
		errorChan := make(chan error, len(sqsEvent.Records))

		svc := sqs.New(session.Must(session.NewSession()))

		for _, sqsMessage := range sqsEvent.Records {
			wg.Add(1)
			go func(msg events.SQSMessage) {
				defer wg.Done()
				err := s.ProcessEvent(msg.Body)
				if err != nil {
					errorChan <- fmt.Errorf("failed to process message %s: %v", msg.MessageId, err)
					// Change visibility timeout for failed messages
					_, changeVisibilityErr := svc.ChangeMessageVisibility(&sqs.ChangeMessageVisibilityInput{
						QueueUrl:          &msg.EventSourceARN,
						ReceiptHandle:     &msg.ReceiptHandle,
						VisibilityTimeout: aws.Int64(0), // Since the message failed, we want it to be immediately visible again
					})

					if changeVisibilityErr != nil {
						log.Printf("Error changing visibility of message with ID %s: %v", msg.MessageId, changeVisibilityErr)
					}
					return
				}
				svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      &msg.EventSourceARN,
					ReceiptHandle: &msg.ReceiptHandle,
				})
				if err != nil {
					log.Printf("Error deleting message with ID %s: %v", msg.MessageId, err)
				}
			}(sqsMessage)
		}

		wg.Wait()
		close(errorChan)

		errorMessages := []string{}

		for err := range errorChan {
			errorMessages = append(errorMessages, err.Error())
		}
		if len(errorMessages) > 0 {
			return fmt.Errorf("failed to process %d messages: %v", len(errorMessages), errorMessages)
		}

		return nil
	}

	lambda.Start(lambdaEventHandler)

	return nil
}
