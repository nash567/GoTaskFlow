package repository

import (
	"fmt"
	"strings"

	"github.com/GoTaskFlow/internal/service/task/model"
)

func buildUpdateTaskFilter(filter *model.UpdateTask) string {
	if filter == nil {
		return ""
	}

	q := `update tasks SET `
	if filter.Title != "" {
		q += fmt.Sprintf("%s = '%s',", "title", filter.Title)

	}
	if filter.Description != "" {
		q += fmt.Sprintf("%s = '%s',", "description", filter.Description)

	}

	if filter.Status != "" {
		q += fmt.Sprintf("%s = '%s',", "status", filter.Status)

	}
	if filter.AssignedTo != "" {
		q += fmt.Sprintf("%s = '%s',", "assigned_to", filter.AssignedTo)

	}

	if filter.Comments != nil {
		q += fmt.Sprintf("%s = '%s',", "comments", filter.Comments)

	}
	if !filter.DueDate.IsZero() {
		q += fmt.Sprintf("%s = '%s',", "due_date", filter.DueDate)

	}

	q = strings.TrimSuffix(q, ",")

	q += fmt.Sprintf("WHERE id = '%s'", filter.ID)
	return q
}
