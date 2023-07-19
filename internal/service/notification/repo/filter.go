package repository

import (
	"fmt"

	"github.com/GoTaskFlow/internal/service/notification/model"
)

func buildCreateNotificationFilter(notifications []model.Notification) (string, []interface{}) {
	query := createNotification
	values := []interface{}{}
	for i, n := range notifications {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		values = append(values, n.Message, n.UserID, n.TaskID, n.Status)
		if i < len(notifications)-1 {
			query += ", "
		}
	}

	return query, values

}
