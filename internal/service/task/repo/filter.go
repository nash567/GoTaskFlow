package repository

import (
	"github.com/GoTaskFlow/internal/service/task/model"
	dbHelper "github.com/GoTaskFlow/pkg/db/helper"
)

const (
	taskStepsTable string = "task_steps"
)

func buildTaskStepFilter(filter *model.Filter) (string, []interface{}) {
	if filter == nil {
		return "", nil
	}
	f := &dbHelper.Filters{}

	if filter.TaskID != nil {
		f.AppendInFilter(taskStepsTable, "task_id", toInterfaceArr(filter.TaskID)...)
	}

	if filter.StepName != nil {
		f.AppendInFilter(taskStepsTable, "name", toInterfaceArr(filter.StepName)...)
	}

	return f.Query(dbHelper.LogicalOperatorAnd, true), f.Params()
}

func toInterfaceArr(v []string) []interface{} {
	out := make([]interface{}, 0, len(v))
	for _, i := range v {
		out = append(out, i)
	}
	return out
}
