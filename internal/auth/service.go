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

	user := User{
		Name:     strings.TrimSpace(name),
		Email:    email,
		Password: string(hash),
	}

	return s.repo.Create(ctx, user)
}

func (s *Service) Login(
	ctx context.Context,
	email, password string,
) (*User, error) {
	u, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password),
	)

	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	return u, nil
}
