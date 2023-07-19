package repository

import (
	"fmt"
	"strings"

	"github.com/GoTaskFlow/internal/service/task/model"
)

func buildUpdateTaskFilter(filter *model.UpdateTask) (string, error) {
	if filter == nil {
		return "", fmt.Errorf("")
	}

	q := `update tasks SET `
	fieldsToupdate := 0
	if filter.Title != "" {
		q += fmt.Sprintf("%s = '%s',", "title", filter.Title)
		fieldsToupdate++

	}
	if filter.Description != "" {
		q += fmt.Sprintf("%s = '%s',", "description", filter.Description)
		fieldsToupdate++
	}

	if filter.Status != "" {
		q += fmt.Sprintf("%s = '%s',", "status", filter.Status)
		fieldsToupdate++
	}
	if filter.AssignedTo != "" {
		q += fmt.Sprintf("%s = '%s',", "assigned_to", filter.AssignedTo)
		fieldsToupdate++
	}

	if filter.Comments != nil {
		q += fmt.Sprintf("%s = '%s',", "comments", filter.Comments)
		fieldsToupdate++
	}
	if !filter.DueDate.IsZero() {
		q += fmt.Sprintf("%s = '%s',", "due_date", filter.DueDate)
		fieldsToupdate++
	}

	if fieldsToupdate <= 0 {
		return "", model.NewNoFieldsToUpdateError()
	}
	q = strings.TrimSuffix(q, ",")

	q += fmt.Sprintf("WHERE id = '%s'", filter.ID)
	return q, nil
}
