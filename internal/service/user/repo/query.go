package repository

const (
	addQuery = `Insert into users (id,name,email,password,active) values ($1,$2,$3,$4,$5) returning id`
	getUsers = `Select * from users`
	// getUserByIDQuery = `Select * from users`
)
