package repository

const (
	getAllQuery        = `select * from notification`
	getByIdQuery       = `select * from notification where id = $1`
	createNotification = `Insert into notification (message,user_id,task_id,status) values `
)
