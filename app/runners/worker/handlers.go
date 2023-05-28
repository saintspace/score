package worker

import (
	"encoding/json"
	"fmt"
	"score/app/models"
)

func (s *Worker) HandleAccountConfirmationTask(eventString string) error {
	event := models.AccountConfirmationTaskPlatformEvent{}
	err := json.Unmarshal([]byte(eventString), &event)
	if err != nil {
		return fmt.Errorf("error parsing event details: %s", err.Error())
	}
	return s.userService.CreateUser(event.EmailAddress, event.UserName)
}

func (s *Worker) HandleEmailBounce(eventString string) error {
	event := models.EmailBouncePlatformEvent{}
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

func (s *Worker) HandleEmailComplaint(eventString string) error {
	event := models.EmailComplaintPlatformEvent{}
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

func (s *Worker) HandleEmailSendTask(eventString string) error {
	event := models.EmailSendTaskPlatformEvent{}
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
