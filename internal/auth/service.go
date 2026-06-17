package auth

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo UserRepository
}

func NewService(repo UserRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(
	ctx context.Context,
	name, email, password string,
) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	name = strings.TrimSpace(name)

	if name == "" {
		return nil, errors.New("name is required")
	}

	user := User{
		Name:     name,
		Email:    email,
		Password: string(hash),
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) Login(
	ctx context.Context,
	email, password string,
) (*User, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return &User{}, errors.New("email is required")
	}

	u, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		return &User{}, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password),
	)

	if err != nil {
		return &User{}, err
	}

	return u, nil
}
