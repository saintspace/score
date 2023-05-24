package handler

type EventHandler struct {
	emailService EmailService
	datastore    Datastore
}

func New(emailService EmailService, datastore Datastore) *EventHandler {
	return &EventHandler{
		emailService: emailService,
		datastore:    datastore,
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

type Datastore interface {
	CreateUser(email, cognitoUserName string) error
}
