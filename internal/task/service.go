package task

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type Service struct {
	repo TasksRepository
}

func NewService(repo TasksRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetTasks(ctx context.Context, user_id uuid.UUID) ([]Task, error) {
	return s.repo.GetTasks(ctx, user_id)
}

func (s *Service) Create(
	ctx context.Context,
	user_id uuid.UUID,
	title, description string, is_completed bool,
) (*Task, error) {
	title = strings.TrimSpace(title)

	if title == "" {
		return nil, errors.New("title is required")
	}

	t := Task{
		Title:       title,
		Description: &description,
		Completed:   is_completed,
		UserID:      user_id,
	}

	return s.repo.Create(ctx, t)
}

func (s *Service) GetTaskById(ctx context.Context, id, user_id uuid.UUID) (*Task, error) {
	t, err := s.repo.FindByID(ctx, id, user_id)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) Update(ctx context.Context, id, user_id uuid.UUID, is_completed *bool, title, description *string) (*Task, error) {
	t, err := s.repo.FindByID(ctx, id, user_id)
	if err != nil {
		return nil, err
	}

	if title != nil {
		t.Title = *title
	}

	if description != nil {
		t.Description = description
	}

	if is_completed != nil {
		t.Completed = *is_completed
	}
	t, err = s.repo.Update(ctx, id, user_id, t)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) Delete(ctx context.Context, id, user_id uuid.UUID) error {
	return s.repo.Delete(ctx, id, user_id)
}
