package task

import (
	"net/http"

	api "github.com/GoTaskFlow/api/http"
	log "github.com/GoTaskFlow/pkg/logger/model"
	"github.com/pborman/uuid"

	"github.com/GoTaskFlow/internal/service/task/model"
	taskModel "github.com/GoTaskFlow/internal/service/task/model"
	"github.com/gorilla/mux"
)

type resource struct {
	service taskModel.Service
	log     log.Logger
}

func RegisterHandlers(router *mux.Router, service model.Service, log log.Logger) {
	res := resource{service, log}

	taskRouter := router.PathPrefix("/task").Subrouter()
	api.RegisterHandlers(taskRouter, http.MethodGet, "", nil, res.get)
	api.RegisterHandlers(taskRouter, http.MethodGet, "/{id}", nil, res.getById)
	api.RegisterHandlers(taskRouter, http.MethodPost, "", nil, res.add)
	api.RegisterHandlers(taskRouter, http.MethodPatch, "/{id}", nil, res.updateTask)

}

func (res *resource) get(w http.ResponseWriter, r *http.Request) {
	tasks, err := res.service.Get(r.Context())
	if err != nil {
		res.log.Error(err.Error())
		api.Write(w, http.StatusInternalServerError, api.NewResponse(false, "failed to get tasks", err))
		return
	}
	res.log.Info("get all tasks")
	api.Write(w, http.StatusOK, tasks)
}

func (res *resource) getById(w http.ResponseWriter, r *http.Request) {
	task, err := res.service.GetTaskByID(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		api.Write(w, http.StatusInternalServerError, api.NewResponse(false, "failed to get task", err))
		return
	}
	api.Write(w, http.StatusOK, task)
}

func (res *resource) add(w http.ResponseWriter, r *http.Request) {
	var task *taskModel.Task
	err := api.SanitizeRequest(r, &task)
	if err != nil {
		api.Write(w, http.StatusInternalServerError, api.NewResponse(false, api.ErrorInvalidRequest, err))
		return

	}
	task.ID = uuid.New()
	err = res.service.Add(r.Context(), task)
	if err != nil {
		api.Write(w, http.StatusInternalServerError, api.NewResponse(false, "failed to add task", err))
		return
	}
	api.Write(w, http.StatusOK, api.NewResponse(true, "task has been added", nil))
}

func (res *resource) updateTask(w http.ResponseWriter, r *http.Request) {
	var task taskModel.UpdateTask
	err := api.SanitizeRequest(r, &task)
	if err != nil {
		api.Write(w, http.StatusInternalServerError, api.NewResponse(false, api.ErrorInvalidRequest, err))
		return
	}
	task.ID = mux.Vars(r)["id"]
	err = res.service.UpdateTask(r.Context(), &task)
	if err != nil {
		api.Write(w, http.StatusInternalServerError, api.NewResponse(false, "failed to get task", err))
		return
	}
	api.Write(w, http.StatusOK, "task has been updated")
}
