package dynamodb

import (
	"fmt"
	"score/app/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDB is responsible for interfacing with AWS DynamoDB Service
type DynamoDB struct {
	svc    *dynamodb.DynamoDB
	config iConfig
}

func New(awsSession *session.Session, config iConfig) *DynamoDB {
	dynamoDBSession := dynamodb.New(awsSession)
	return &DynamoDB{
		svc:    dynamoDBSession,
		config: config,
	}
}

type iConfig interface {
	EmailSubscriptionsTableName() string
}

func (s *DynamoDB) itemExists(
	tableName string,
	key map[string]*dynamodb.AttributeValue,
) (bool, error) {
	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}
	result, err := s.svc.GetItem(getInput)
	if err != nil {
		return false, err
	}
	if result.Item != nil {
		return true, nil
	}
	return false, nil
}

func (s *DynamoDB) updateItem(
	tableName string,
	key map[string]*dynamodb.AttributeValue,
	expressionAttributeValues map[string]*dynamodb.AttributeValue,
	updateExpression string,
	returnValues string,
) error {
	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		ExpressionAttributeValues: expressionAttributeValues,
		UpdateExpression:          aws.String(updateExpression),
		ReturnValues:              aws.String(returnValues),
	}
	_, err := s.svc.UpdateItem(updateInput)
	return err
}

func (s *DynamoDB) EmailSubscriptionItemExists(email string) (bool, error) {
	tableName := s.config.EmailSubscriptionsTableName()
	key := map[string]*dynamodb.AttributeValue{
		"email": {
			S: aws.String(email),
		},
	}
	return s.itemExists(tableName, key)
}

func (s *DynamoDB) AddComplaintToEmailSubscription(email, complaintDetails string, complaintDateUnix int64) error {
	tableName := s.config.EmailSubscriptionsTableName()
	return s.updateItem(
		tableName,
		map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		map[string]*dynamodb.AttributeValue{
			":hasComplaint": {
				BOOL: aws.Bool(true),
			},
			":complaintDate": {
				N: aws.String(fmt.Sprintf("%d", complaintDateUnix)),
			},
			":complaintDetails": {
				S: aws.String(complaintDetails),
			},
		},
		"set has_complaint = :hasComplaint, complaint_date = :complaintDate, complaint_details = :complaintDetails",
		"NONE",
	)
}

func (s *DynamoDB) AddBounceToEmailSubscription(
	email,
	bounceType string,
	bounceDetails string,
	bounceDateUnix int64,
) error {
	tableName := s.config.EmailSubscriptionsTableName()
	return s.updateItem(
		tableName,
		map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		map[string]*dynamodb.AttributeValue{
			":hasBounce": {
				BOOL: aws.Bool(true),
			},
			":bounceDate": {
				N: aws.String(fmt.Sprintf("%d", bounceDateUnix)),
			},
			":bounceType": {
				S: aws.String(bounceType),
			},
			":bounceDetails": {
				S: aws.String(bounceDetails),
			},
		},
		"set has_bounce = :hasBounce, bounce_date = :bounceDate, bounce_type = :bounceType, bounce_details = :bounceDetails",
		"NONE",
	)
}

func (s *DynamoDB) GetEmailSubscription(email string) (*models.EmailSubscription, error) {
	tableName := s.config.EmailSubscriptionsTableName()
	getInput := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	}
	result, err := s.svc.GetItem(getInput)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}
	emailSubscription := &models.EmailSubscription{}
	err = dynamodbattribute.UnmarshalMap(result.Item, emailSubscription)
	if err != nil {
		return nil, err
	}
	return emailSubscription, nil
}
