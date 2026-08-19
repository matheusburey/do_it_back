package task

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

type fakeTaskRepo struct {
	createFn   func(ctx context.Context, t Task) (*Task, error)
	getTasksFn func(ctx context.Context, user_id uuid.UUID) ([]Task, error)
	findByIDFn func(ctx context.Context, id, user_id uuid.UUID) (*Task, error)
	updateFn   func(ctx context.Context, id, user_id uuid.UUID, task *Task) (*Task, error)
	deleteFn   func(ctx context.Context, id, user_id uuid.UUID) error
}

func (f *fakeTaskRepo) GetTasks(ctx context.Context, user_id uuid.UUID) ([]Task, error) {
	return f.getTasksFn(ctx, user_id)
}

func (f *fakeTaskRepo) Create(ctx context.Context, t Task) (*Task, error) {
	return f.createFn(ctx, t)
}

func (f *fakeTaskRepo) FindByID(ctx context.Context, id, user_id uuid.UUID) (*Task, error) {
	return f.findByIDFn(ctx, id, user_id)
}

func (f *fakeTaskRepo) Update(ctx context.Context, id, user_id uuid.UUID, task *Task) (*Task, error) {
	return f.updateFn(ctx, id, user_id, task)
}

func (f *fakeTaskRepo) Delete(ctx context.Context, id, user_id uuid.UUID) error {
	return f.deleteFn(ctx, id, user_id)
}

func TestCreate(t *testing.T) {
	user_id := uuid.New()
	ts := NewService(&fakeTaskRepo{
		createFn: func(ctx context.Context, t Task) (*Task, error) {
			return &t, nil
		},
	})

	task_title := "Title Task"
	task_description := "Description"

	task, err := ts.Create(context.Background(), user_id, task_title, task_description, false)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task == nil {
		t.Fatalf("expected task, got nil")
	}

	if task.Title != task_title {
		t.Errorf("expected title to be %s, got %s", task_title, task.Title)
	}
}

func TestCreate_TrimsTitle(t *testing.T) {
	ts := NewService(&fakeTaskRepo{
		createFn: func(ctx context.Context, t Task) (*Task, error) {
			return &t, nil
		},
	})

	task, err := ts.Create(context.Background(), uuid.New(), "  Title Task  ", "Description", false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.Title != "Title Task" {
		t.Errorf("expected trimmed title %q, got %q", "Title Task", task.Title)
	}
}

func TestCreate_EmptyTitle(t *testing.T) {
	called := false
	ts := NewService(&fakeTaskRepo{
		createFn: func(ctx context.Context, t Task) (*Task, error) {
			called = true
			return &t, nil
		},
	})

	_, err := ts.Create(context.Background(), uuid.New(), "   ", "Description", false)

	if !errors.Is(err, ErrTitleRequired) {
		t.Fatalf("expected ErrTitleRequired, got %v", err)
	}

	if called {
		t.Errorf("expected repository.Create not to be called when title is blank")
	}
}

func TestGetTaskById_Success(t *testing.T) {
	user_id := uuid.New()
	task_id := uuid.New()
	want := &Task{ID: task_id, UserID: user_id, Title: "Task 1"}

	ts := NewService(&fakeTaskRepo{
		findByIDFn: func(ctx context.Context, id, uid uuid.UUID) (*Task, error) {
			return want, nil
		},
	})

	got, err := ts.GetTaskById(context.Background(), task_id, user_id)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got != want {
		t.Errorf("expected task %v, got %v", want, got)
	}
}

func TestGetTaskById_NotFound(t *testing.T) {
	ts := NewService(&fakeTaskRepo{
		findByIDFn: func(ctx context.Context, id, uid uuid.UUID) (*Task, error) {
			return nil, ErrTaskNotFound
		},
	})

	_, err := ts.GetTaskById(context.Background(), uuid.New(), uuid.New())

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestUpdate_PartialUpdate(t *testing.T) {
	user_id := uuid.New()
	task_id := uuid.New()
	original_description := "original description"
	existing := &Task{
		ID:          task_id,
		UserID:      user_id,
		Title:       "Original Title",
		Description: &original_description,
		Completed:   false,
	}

	ts := NewService(&fakeTaskRepo{
		findByIDFn: func(ctx context.Context, id, uid uuid.UUID) (*Task, error) {
			return existing, nil
		},
		updateFn: func(ctx context.Context, id, uid uuid.UUID, task *Task) (*Task, error) {
			return task, nil
		},
	})

	is_completed := true
	got, err := ts.Update(context.Background(), task_id, user_id, &is_completed, nil, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.Title != "Original Title" {
		t.Errorf("expected title to stay %q, got %q", "Original Title", got.Title)
	}

	if !got.Completed {
		t.Errorf("expected completed to be true")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	ts := NewService(&fakeTaskRepo{
		findByIDFn: func(ctx context.Context, id, uid uuid.UUID) (*Task, error) {
			return nil, ErrTaskNotFound
		},
	})

	_, err := ts.Update(context.Background(), uuid.New(), uuid.New(), nil, nil, nil)

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestDelete_Success(t *testing.T) {
	ts := NewService(&fakeTaskRepo{
		deleteFn: func(ctx context.Context, id, uid uuid.UUID) error {
			return nil
		},
	})

	err := ts.Delete(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	ts := NewService(&fakeTaskRepo{
		deleteFn: func(ctx context.Context, id, uid uuid.UUID) error {
			return ErrTaskNotFound
		},
	})

	err := ts.Delete(context.Background(), uuid.New(), uuid.New())

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}
