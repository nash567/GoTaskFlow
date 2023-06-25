package repository

const (
	addQuery     = `Insert into users (id,name,email,password,active) values ($1,$2,$3,$4,$5) returning id`
	getAllQuery  = `Select * from users`
	getByIdQuery = `Select * from users where id=$1`
)
