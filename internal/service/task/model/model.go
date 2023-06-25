package model

import "github.com/lib/pq"

type Task struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	AssignedBy  string         `json:"assigned_by"`
	AssignedTo  string         `json:"assigned_to"`
	Comments    pq.StringArray `json:"comments"`
	DueDate     string         `json:"due_date"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
}
