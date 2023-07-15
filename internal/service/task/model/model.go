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

type TaskWorkflowInput struct {
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

type TaskStep struct {
	ID        *int
	Name      *string
	StepOrder *int
	TaskID    *string
	Status    *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
