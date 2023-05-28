package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"score/app/config"
	"score/app/logger"
	"score/app/runners/server"
	"score/app/runners/worker"
	"score/app/services/aws/dynamodb"
	"score/app/services/aws/ses"
	"score/app/services/aws/sns"
	"score/app/services/datastore"
	"score/app/services/email"
	"score/app/services/eventpub"
	"score/app/services/mysql"
	"score/app/services/user"

	"github.com/aws/aws-sdk-go/aws/session"
)

//go:embed dist/*
var dist embed.FS

type ExecutionMode string

const (
	ServerMode    ExecutionMode = "server"
	WorkerMode    ExecutionMode = "worker"
	ConfirmerMode ExecutionMode = "confirmer"
	PreauthMode   ExecutionMode = "preauth"
)

func main() {
	executionMode := getExecutionMode()
	appRunner, err := buildAppRunner(executionMode)
	if err != nil {
		panic(err)
	}
	err = appRunner.Run()
	if err != nil {
		panic(err)
	}
}

type Runner interface {
	Run() error
}

// Checks the 'SCORE_MODE' environment variable to determine the execution mode. If it is not set, it will check the 'mode' command line option. If that is not set, it will default to 'server' mode.
func getExecutionMode() ExecutionMode {
	if os.Getenv("SCORE_MODE") != "" {
		return ExecutionMode(os.Getenv("SCORE_MODE"))
	}
	executionModePtr := flag.String("mode", "server", "The execution mode of the application.")
	flag.Parse()
	return ExecutionMode(*executionModePtr)
}

func buildAppRunner(mode ExecutionMode) (Runner, error) {

	awsProfile := os.Getenv("SCORE_AWS_PROFILE")
	awsSessionOptions := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}
	if awsProfile != "" {
		awsSessionOptions.Profile = awsProfile
	}
	awsSession := session.Must(session.NewSessionWithOptions(awsSessionOptions))
	configService := config.New(awsSession)
	err := configService.InitializeParameters()
	if err != nil {
		return nil, fmt.Errorf("error initializing app config: %v", err.Error())
	}

	// Set up the logger
	loggerService, err := logger.New(false)
	if err != nil {
		return nil, fmt.Errorf("error initializing logger: %v", err.Error())
	}

	sesService := ses.New(awsSession)
	snsService := sns.New(awsSession, configService)
	mysqlService := mysql.New(configService)
	dynamoDbService := dynamodb.New(awsSession, configService)
	datastoreService := datastore.New(dynamoDbService, mysqlService)
	eventPublisherService := eventpub.New(snsService, configService)
	userService := user.New(datastoreService)
	emailService := email.New(sesService, datastoreService, loggerService, eventPublisherService)

	switch mode {
	case ServerMode:
		return server.New(emailService, loggerService, dist), nil
	case WorkerMode:
		return worker.New(emailService, userService, loggerService), nil
	case ConfirmerMode:
		return worker.New(emailService, userService, loggerService), nil
	case PreauthMode:
		return worker.New(emailService, userService, loggerService), nil
	default:
		return nil, fmt.Errorf("invalid execution mode: %s", mode)
	}
}
