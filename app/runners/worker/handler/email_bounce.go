package handler

import (
	"encoding/json"
	"fmt"
	"time"
)

func (s *EventHandler) EmailBounce(eventString string) error {
	event := EmailBounceEvent{}
	err := json.Unmarshal([]byte(eventString), &event)
	if err != nil {
		return fmt.Errorf("error parsing event details: %s", err.Error())
	}
	emailAddresses := []string{}
	for _, recipient := range event.Bounce.BouncedRecipients {
		emailAddresses = append(emailAddresses, recipient.EmailAddress)
	}
	return s.emailService.ProcessEmailBounce(
		emailAddresses,
		event.Bounce.BounceType,
		event.Bounce.BounceSubType,
		eventString,
		event.Bounce.Timestamp.Unix(),
	)
}

type EmailBounceEvent struct {
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
