package handler

import (
	"encoding/json"
	"fmt"
	"time"
)

func (s *EventHandler) EmailComplaint(eventString string) error {
	event := EmailComplaintEvent{}
	err := json.Unmarshal([]byte(eventString), &event)
	if err != nil {
		return fmt.Errorf("error parsing event details: %s", err.Error())
	}
	emailAddresses := []string{}
	for _, recipient := range event.Complaint.ComplainedRecipients {
		emailAddresses = append(emailAddresses, recipient.EmailAddress)
	}
	return s.emailService.ProcessEmailComplaint(emailAddresses, eventString, event.Complaint.Timestamp.Unix())
}

type EmailComplaintEvent struct {
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
