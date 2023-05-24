package handler

import (
	"encoding/json"
	"fmt"
)

type AccountConfirmationTaskEvent struct {
	UserName     string `json:"userName"`
	EmailAddress string `json:"emailAddress"`
}

func (s *EventHandler) AccountConfirmationTask(eventString string) error {
	event := AccountConfirmationTaskEvent{}
	err := json.Unmarshal([]byte(eventString), &event)
	if err != nil {
		return fmt.Errorf("error parsing event details: %s", err.Error())
	}
	return s.datastore.CreateUser(event.EmailAddress, event.UserName)
}
