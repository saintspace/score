package models

import (
	"time"
)

type PlatformEvent struct {
	EventName     PlatformEventName `json:"eventType"`
	CorrelationId string            `json:"correlationId"`
	EventDetails  string            `json:"eventDetails"`
}

type PlatformEventName string
type platformEventNames struct {
	EmailBounce             PlatformEventName
	EmailComplaint          PlatformEventName
	EmailSendTask           PlatformEventName
	AccountConfirmationTask PlatformEventName
}

var PlatformEventNames = platformEventNames{
	EmailBounce:             "Bounce",
	EmailComplaint:          "Complaint",
	EmailSendTask:           "email-send-task",
	AccountConfirmationTask: "account-confirmation-task",
}

type EmailSendTaskPlatformEvent struct {
	TemplateName  string            `json:"templateName"`
	SenderAddress string            `json:"senderAddress"`
	SubjectLine   string            `json:"subjectLine"`
	ToAddresses   []string          `json:"toAddresses"`
	Parameters    map[string]string `json:"parameters"`
}

type EmailComplaintPlatformEvent struct {
	EventType string         `json:"eventType"`
	Complaint EmailComplaint `json:"complaint"`
	Mail      ComplaintMail  `json:"mail"`
}

type EmailComplaint struct {
	FeedbackID            string                     `json:"feedbackId"`
	ComplaintSubType      *string                    `json:"complaintSubType"`
	ComplainedRecipients  []EmailComplainedRecipient `json:"complainedRecipients"`
	Timestamp             time.Time                  `json:"timestamp"`
	UserAgent             string                     `json:"userAgent"`
	ComplaintFeedbackType string                     `json:"complaintFeedbackType"`
	ArrivalDate           time.Time                  `json:"arrivalDate"`
}

type EmailComplainedRecipient struct {
	EmailAddress string `json:"emailAddress"`
}

type ComplaintMail struct {
	Timestamp        time.Time                   `json:"timestamp"`
	Source           string                      `json:"source"`
	SourceARN        string                      `json:"sourceArn"`
	SendingAccountID string                      `json:"sendingAccountId"`
	MessageID        string                      `json:"messageId"`
	Destination      []string                    `json:"destination"`
	HeadersTruncated bool                        `json:"headersTruncated"`
	Headers          []EmailComplaintHeader      `json:"headers"`
	CommonHeaders    EmailComplaintCommonHeaders `json:"commonHeaders"`
	Tags             map[string][]string         `json:"tags"`
}

type EmailComplaintHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EmailComplaintCommonHeaders struct {
	From      []string `json:"from"`
	To        []string `json:"to"`
	MessageID string   `json:"messageId"`
	Subject   string   `json:"subject"`
}

type EmailBouncePlatformEvent struct {
	EventType string      `json:"eventType"`
	Bounce    EmailBounce `json:"bounce"`
	Mail      BounceMail  `json:"mail"`
}

type EmailBounce struct {
	FeedbackID        string                  `json:"feedbackId"`
	BounceType        string                  `json:"bounceType"`
	BounceSubType     string                  `json:"bounceSubType"`
	BouncedRecipients []EmailBouncedRecipient `json:"bouncedRecipients"`
	Timestamp         time.Time               `json:"timestamp"`
	ReportingMTA      string                  `json:"reportingMTA"`
}

type EmailBouncedRecipient struct {
	EmailAddress   string `json:"emailAddress"`
	Action         string `json:"action"`
	Status         string `json:"status"`
	DiagnosticCode string `json:"diagnosticCode"`
}

type BounceMail struct {
	Timestamp        time.Time                `json:"timestamp"`
	Source           string                   `json:"source"`
	SourceARN        string                   `json:"sourceArn"`
	SendingAccountID string                   `json:"sendingAccountId"`
	MessageID        string                   `json:"messageId"`
	Destination      []string                 `json:"destination"`
	HeadersTruncated bool                     `json:"headersTruncated"`
	Headers          []EmailBounceHeader      `json:"headers"`
	CommonHeaders    EmailBounceCommonHeaders `json:"commonHeaders"`
	Tags             map[string][]string      `json:"tags"`
}

type EmailBounceHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type EmailBounceCommonHeaders struct {
	From      []string `json:"from"`
	To        []string `json:"to"`
	MessageID string   `json:"messageId"`
	Subject   string   `json:"subject"`
}

type AccountConfirmationTaskPlatformEvent struct {
	UserName     string `json:"userName"`
	EmailAddress string `json:"emailAddress"`
}
