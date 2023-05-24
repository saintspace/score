package datastore

import (
	"score/app/models"
)

type Datastore struct {
	kvStore      KeyValueStore
	relationalDB RelationalDB
}

func New(kvStore KeyValueStore, relationalDB RelationalDB) *Datastore {
	return &Datastore{
		kvStore:      kvStore,
		relationalDB: relationalDB,
	}
}

type KeyValueStore interface {
	EmailSubscriptionItemExists(email string) (bool, error)
	AddComplaintToEmailSubscription(email, complaintDetails string, complaintDateUnix int64) error
	AddBounceToEmailSubscription(email, bounceType, bounceDetails string, bounceDateUnix int64) error
	GetEmailSubscription(email string) (*models.EmailSubscription, error)
}

type RelationalDB interface {
	CreateUser(email, cognitoUserName string) error
}

func (s *Datastore) EmailSubscriptionItemExists(email string) (bool, error) {
	return s.kvStore.EmailSubscriptionItemExists(email)
}

func (s *Datastore) AddComplaintToEmailSubscription(email, complaintDetails string, complaintDateUnix int64) error {
	return s.kvStore.AddComplaintToEmailSubscription(email, complaintDetails, complaintDateUnix)
}

func (s *Datastore) AddBounceToEmailSubscription(email, bounceType, bounceDetails string, bounceDateUnix int64) error {
	return s.kvStore.AddBounceToEmailSubscription(email, bounceType, bounceDetails, bounceDateUnix)
}

func (s *Datastore) GetEmailSubscription(email string) (*models.EmailSubscription, error) {
	return s.kvStore.GetEmailSubscription(email)
}

func (s *Datastore) CreateUser(email, cognitoUserName string) error {
	return s.relationalDB.CreateUser(email, cognitoUserName)
}
