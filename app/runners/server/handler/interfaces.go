package handler

type RouteHandler struct {
	emailService  EmailService
	loggerService LoggerService
}

func New(emailService EmailService, loggerService LoggerService) *RouteHandler {
	return &RouteHandler{
		emailService:  emailService,
		loggerService: loggerService,
	}
}

type EmailService interface {
	CreateEmailSubscription(email string) error
	EmailSubscriptionExists(email string) (bool, error)
	IsValidEmail(email string) bool
	VerifyEmailWithSubscriptionToken(token string) error
}

type LoggerService interface {
	InfoWithContext(message string, keysAndValues ...interface{})
	ErrorWithContext(message string, keysAndValues ...interface{})
	DebugWithContext(message string, keysAndValues ...interface{})
}
