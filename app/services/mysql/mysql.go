package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type MySQL struct {
	config Config
}

func New(config Config) *MySQL {
	return &MySQL{
		config: config,
	}
}

type Config interface {
	RelationalDatabaseConnectionString() string
}

func (s *MySQL) CreateUser(email, cognitoUserName string) error {
	db, err := sql.Open("mysql", s.config.RelationalDatabaseConnectionString())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO users (id, email, cognito_user_id) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	id := uuid.New().String()

	_, err = stmt.Exec(id, email, cognitoUserName)
	if err != nil {
		return fmt.Errorf("failed to execute SQL statement: %v", err)
	}

	return nil
}
