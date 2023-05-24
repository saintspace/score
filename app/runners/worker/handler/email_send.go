package handler

import (
	"encoding/json"
	"fmt"
)

type EmailSendTaskEvent struct {
	TemplateName  string            `json:"templateName"`
	SenderAddress string            `json:"senderAddress"`
	SubjectLine   string            `json:"subjectLine"`
	ToAddresses   []string          `json:"toAddresses"`
	Parameters    map[string]string `json:"parameters"`
}

func (s *EventHandler) EmailSendTask(eventString string) error {
	event := EmailSendTaskEvent{}
	err := json.Unmarshal([]byte(eventString), &event)
	if err != nil {
		return fmt.Errorf("error parsing event details: %s", err.Error())
	}
	err = s.emailService.SendTemplatedEmail(
		event.TemplateName,
		event.Parameters,
		event.SubjectLine,
		event.SenderAddress,
		event.ToAddresses,
	)
	if err != nil {
		return fmt.Errorf("error sending templated email: %s", err.Error())
	}
	return nil
}
