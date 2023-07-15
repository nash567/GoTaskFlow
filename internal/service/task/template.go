package task

import (
	"bytes"
	"html/template"
)

var createTaskMail = `<h1>Task has been created for you</h1>`
var notifyDueDate = `<h1>Task has reached the due date! Please complete it</h1>`

func (s *Service) getBody(body string) string {
	tmpl, err := template.New("taskCreation").Parse(body)
	if err != nil {
		return ""
	}
	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, ""); err != nil {
		panic(err)
	}
	return tpl.String()
}
