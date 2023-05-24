package user

type UserService struct {
	datastore Datastore
}

func New(datastore Datastore) *UserService {
	return &UserService{
		datastore: datastore,
	}
}

type Datastore interface {
	CreateUser(email, cognitoUserName string) error
}

func (s *UserService) CreateUser(email, cognitoUserName string) error {
	return s.datastore.CreateUser(email, cognitoUserName)
}
