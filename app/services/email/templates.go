package email

import (
	"bytes"
	"fmt"
	"text/template"
)

func generateEmailBodyFromTemplate(templateName string, params map[string]string) (string, error) {
	templates := EmailTemplateCollection
	emailTemplates := map[string]*template.Template{}
	for k, v := range templates {
		emailTemplates[k] = template.Must(template.New("index").Parse(v))
	}
	currentTemplate, ok := emailTemplates[templateName]
	if !ok {
		return "", fmt.Errorf("unknown email template  => %s", templateName)
	}
	emailBody := &bytes.Buffer{}
	err := currentTemplate.Execute(emailBody, params)
	return emailBody.String(), err
}

var EmailTemplateCollection = map[string]string{
	"email-subscription-verification": "Thank you for your interest in SaintSpace!\n\nPlease click on the link below to confirm your subscription. It helps us make sure you're human.\n\n {{.verificationLink}}\n\n",
}
