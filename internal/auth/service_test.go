package auth

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepo struct {
	createFn func(ctx context.Context, u User) (*User, error)
	findFn   func(ctx context.Context, email string) (*User, error)
}

func (f *fakeUserRepo) Create(ctx context.Context, u User) (*User, error) {
	return f.createFn(ctx, u)
}

func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	return f.findFn(ctx, email)
}

func TestRegister(t *testing.T) {
	as := NewService(&fakeUserRepo{
		createFn: func(ctx context.Context, u User) (*User, error) {
			return &u, nil
		},
	})

	user, err := as.Register(context.Background(), "John Doe", "8tH9p@example.com", "password")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user == nil {
		t.Fatalf("expected user, got nil")
	}

	if user.Name != "John Doe" {
		t.Errorf("expected name to be %s, got %s", "John Doe", user.Name)
	}
}

func TestRegister_HashesPassword(t *testing.T) {
	as := NewService(&fakeUserRepo{
		createFn: func(ctx context.Context, u User) (*User, error) {
			return &u, nil
		},
	})

	password := "password"
	user, err := as.Register(context.Background(), "John Doe", "8tH9p@example.com", password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Password == password {
		t.Fatalf("expected password to be hashed, got plaintext")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		t.Errorf("expected stored hash to match original password: %v", err)
	}
}

func TestRegister_TrimsName(t *testing.T) {
	as := NewService(&fakeUserRepo{
		createFn: func(ctx context.Context, u User) (*User, error) {
			return &u, nil
		},
	})

	user, err := as.Register(context.Background(), "  John Doe  ", "8tH9p@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Name != "John Doe" {
		t.Errorf("expected trimmed name %q, got %q", "John Doe", user.Name)
	}
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	as := NewService(&fakeUserRepo{
		createFn: func(ctx context.Context, u User) (*User, error) {
			return nil, ErrEmailAlreadyExists
		},
	})

	_, err := as.Register(context.Background(), "John Doe", "8tH9p@example.com", "password")

	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func newHashedUser(t *testing.T, email, password string) *User {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	return &User{
		Email:    email,
		Password: string(hash),
	}
}

func TestLogin_Success(t *testing.T) {
	email := "8tH9p@example.com"
	password := "password"
	stored := newHashedUser(t, email, password)

	as := NewService(&fakeUserRepo{
		findFn: func(ctx context.Context, e string) (*User, error) {
			return stored, nil
		},
	})

	user, err := as.Login(context.Background(), email, password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != email {
		t.Errorf("expected email %q, got %q", email, user.Email)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	email := "8tH9p@example.com"
	stored := newHashedUser(t, email, "password")

	as := NewService(&fakeUserRepo{
		findFn: func(ctx context.Context, e string) (*User, error) {
			return stored, nil
		},
	})

	_, err := as.Login(context.Background(), email, "wrong-password")

	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	as := NewService(&fakeUserRepo{
		findFn: func(ctx context.Context, e string) (*User, error) {
			return nil, ErrUserNotFound
		},
	})

	_, err := as.Login(context.Background(), "missing@example.com", "password")

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}
