package repository

import (
	"github.com/GoTaskFlow/internal/service/user/model"
	dbHelper "github.com/GoTaskFlow/pkg/db/helper"
)

const (
	userTable string = "users"
)

func buildFilter(filter *model.Filter) (string, []interface{}) {
	if filter == nil {
		return "", nil
	}
	f := &dbHelper.Filters{}

	if filter.ID != nil {
		f.AppendInFilter(userTable, "id", toInterfaceArr(filter.ID)...)
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
