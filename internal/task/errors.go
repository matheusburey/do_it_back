package task

import "errors"

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrTitleRequired = errors.New("title is required")
)
