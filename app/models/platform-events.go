package models

type PlatformEvent struct {
	EventName     string `json:"eventType"`
	CorrelationId string `json:"correlationId"`
	EventDetails  string `json:"eventDetails"`
}

type EmailSendTaskPlatformEvent struct {
	TemplateName  string            `json:"templateName"`
	SenderAddress string            `json:"senderAddress"`
	SubjectLine   string            `json:"subjectLine"`
	ToAddresses   []string          `json:"toAddresses"`
	Parameters    map[string]string `json:"parameters"`
}
