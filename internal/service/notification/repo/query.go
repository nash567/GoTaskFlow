package repository

const (
	getAllQuery  = `select * from notification`
	getByIdQuery = `select * from notification where id = $1`
	addQuery     = `Insert into notification (message,user_id,task_id,status) values ($1,$2,$3,$4)`
)
