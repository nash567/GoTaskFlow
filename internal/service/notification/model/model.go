package model

type Notification struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	UserID  string `json:"user_id"`
	TaskID  string `json:"task_id"`
	Status  Status `json:"status"`
}
