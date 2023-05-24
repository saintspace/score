package worker

import (
	"context"
	"fmt"
	"log"
	"score/app/config"
	"score/app/logger"
	"score/app/runners/worker/handler"
	"score/app/services/aws/dynamodb"
	"score/app/services/aws/ses"
	"score/app/services/datastore"
	"score/app/services/email"
	"score/app/services/mysql"
	"score/app/services/user"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var eventRouter *EventRouter

var Run = func() error {
	fmt.Println("Running score in worker mode...")
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	configService := config.New(awsSession)
	err := configService.InitializeParameters()
	if err != nil {
		return fmt.Errorf("error initializing app config: %v", err.Error())
	}

	// Set up the logger
	loggerService, err := logger.New(false)
	if err != nil {
		return fmt.Errorf("error initializing logger: %v", err.Error())
	}

	// Build application dependencies
	sesService := ses.New(awsSession)
	mysqlService := mysql.New(configService)
	dynamoDbService := dynamodb.New(awsSession, configService)
	datastoreService := datastore.New(dynamoDbService, mysqlService)
	userService := user.New(datastoreService)
	emailService := email.New(sesService, datastoreService, loggerService)
	platformEventHandler := handler.New(emailService, userService)
	eventRouter = NewRouter(platformEventHandler, loggerService)

	lambda.Start(lambdaEventHandler)

	return nil
}

func lambdaEventHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(sqsEvent.Records))

	svc := sqs.New(session.Must(session.NewSession()))

	for _, sqsMessage := range sqsEvent.Records {
		wg.Add(1)
		go func(msg events.SQSMessage) {
			defer wg.Done()
			err := eventRouter.ProcessEvent(msg.Body)
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
