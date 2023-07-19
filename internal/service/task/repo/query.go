package repository

const (
	addQuery            = `Insert into tasks (id,title,description,status,assigned_by,assigned_to,comments,due_date)values($1,$2,$3,$4,$5,$6,$7,$8)`
	getQuery            = `Select * from tasks`
	getByIdQuery        = `Select * from tasks where id=$1`
	createTaskStepQuery = `Insert into task_steps(name,status,step_order,task_id) values($1,$2,$3,$4)`

	updateTaskStepQuery = `Update task_steps SET status= $1,updated_at=CURRENT_TIMESTAMP where task_id=$2 and name=$3`
	getTasksWithDueDate = `Select * from tasks WHERE due_date >= CURRENT_DATE AND due_date < CURRENT_DATE + INTERVAL '1 day'`
)
