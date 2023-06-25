package repository

const (
	addQuery     = `Insert into tasks (id,title,description,status,assigned_by,assigned_to,comments,due_date)values($1,$2,$3,$4,$5,$6,$7,$8)`
	getQuery     = `Select * from tasks`
	getByIdQuery = `Select * from tasks where id=$1`
)
