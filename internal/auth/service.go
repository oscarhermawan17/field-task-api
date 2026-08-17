package auth

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	EmployeeID string `json:"employeeId"`
}

type userRecord struct {
	User
	passwordHash []byte
}

// Service owns the temporary in-memory user store for the authentication task.
// Replace it with a repository when the database task begins.
type Service struct {
	usersByEmail map[string]userRecord
	usersByID    map[string]userRecord
}

func NewService(seedPassword string) (*Service, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	seedUser := userRecord{
		User: User{
			ID:         "user-1",
			Name:       "Field Worker",
			Email:      "worker@fieldtask.com",
			EmployeeID: "EMP-001",
		},
		passwordHash: passwordHash,
	}

	return &Service{
		usersByEmail: map[string]userRecord{seedUser.Email: seedUser},
		usersByID:    map[string]userRecord{seedUser.ID: seedUser},
	}, nil
}

func (s *Service) Authenticate(email, password string) (User, error) {
	user, ok := s.usersByEmail[strings.ToLower(strings.TrimSpace(email))]
	if !ok || bcrypt.CompareHashAndPassword(user.passwordHash, []byte(password)) != nil {
		return User{}, ErrInvalidCredentials
	}

	return user.User, nil
}

func (s *Service) UserByID(id string) (User, bool) {
	user, ok := s.usersByID[id]
	if !ok {
		return User{}, false
	}

	return user.User, true
}
