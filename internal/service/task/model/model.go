package model

import (
	"time"

	"github.com/lib/pq"
)

type Task struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	AssignedBy  string         `json:"assigned_by"`
	AssignedTo  string         `json:"assigned_to"`
	Comments    pq.StringArray `json:"comments"`
	DueDate     time.Time      `json:"due_date"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}

type UpdateTask struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	AssignedBy  string         `json:"assigned_by"`
	AssignedTo  string         `json:"assigned_to"`
	Comments    pq.StringArray `json:"comments"`
	DueDate     time.Time      `json:"due_date"`
}

type TaskStep struct {
	ID        *int
	Name      *string
	StepOrder *int
	TaskID    *string
	Status    *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

func (u *UpdateTask) PrepareUpateTask(originalTask *Task) *UpdateTask {

	if u.Title != "" && u.Title == originalTask.Title {
		u.Title = ""
	}
	if u.Description != "" && u.Description == originalTask.Description {
		u.Description = ""
	}
	if u.Status != "" && u.Status == originalTask.Status {
		u.Status = ""
	}
	if !u.DueDate.IsZero() && u.DueDate.Equal(originalTask.DueDate) {
		u.DueDate = time.Time{}
	}

	if originalTask.AssignedTo == u.AssignedTo {
		u.AssignedTo = ""
	}

	return u
}
