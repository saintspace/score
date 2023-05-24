package models

type EmailSubscription struct {
	Email             string `json:"email"`
	Verified          bool   `json:"email_verified"`
	HasComplaint      bool   `json:"has_complaint"`
	HasBounce         bool   `json:"has_bounce"`
	ComplaintDetails  string `json:"complaint_details"`
	ComplaintDateUnix int64  `json:"complaint_date"`
	BounceType        string `json:"bounce_type"`
	BounceDetails     string `json:"bounce_details"`
	BounceDateUnix    int64  `json:"bounce_date"`
	CreationDate      int64  `json:"creation_date"`
	SubscriptionToken string `json:"subscription_token"`
}
