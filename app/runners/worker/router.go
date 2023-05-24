package worker

import (
	"encoding/json"
	"fmt"
)

type EventRouter struct {
	eventHandler EventHandler
	logger       Logger
}

func NewRouter(eventHandler EventHandler, logger Logger) *EventRouter {
	return &EventRouter{
		eventHandler: eventHandler,
		logger:       logger,
	}
}

type EventHandler interface {
	EmailSendTask(eventDetails string) error
	AccountConfirmationTask(eventDetails string) error
	EmailBounce(eventDetails string) error
	EmailComplaint(eventDetails string) error
}

type Logger interface {
	InfoWithContext(message string, keysAndValues ...interface{})
	ErrorWithContext(message string, keysAndValues ...interface{})
	DebugWithContext(message string, keysAndValues ...interface{})
}

type Event struct {
	EventName     string `json:"eventType"`
	CorrelationId string `json:"correlationId"`
	EventDetails  string `json:"eventDetails"`
}

func (s *EventRouter) ProcessEvent(eventString string) error {
	event := Event{}
	err := json.Unmarshal([]byte(eventString), &event)
	if err != nil {
		return fmt.Errorf("error parsing event: %s", err.Error())
	}
	if event.EventName == "Bounce" || event.EventName == "Complaint" {
		event.EventDetails = eventString
	}
	switch event.EventName {
	case "email-send-task":
		err = s.eventHandler.EmailSendTask(event.EventDetails)
	case "account-confirmation-task":
		err = s.eventHandler.AccountConfirmationTask(event.EventDetails)
	case "Bounce":
		err = s.eventHandler.EmailBounce(event.EventDetails)
	case "Complaint":
		err = s.eventHandler.EmailComplaint(event.EventDetails)
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
