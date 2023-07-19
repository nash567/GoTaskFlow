package user

import (
	"net/http"

	api "github.com/GoTaskFlow/api/http"
	"github.com/GoTaskFlow/internal/service/user/model"
	log "github.com/GoTaskFlow/pkg/logger/model"
	"github.com/gorilla/mux"
	"github.com/pborman/uuid"
)

type resource struct {
	service model.Service
	log     log.Logger
}

func RegisterHandlers(router *mux.Router, service model.Service, log log.Logger) {
	res := resource{service, log}

	userRouter := router.PathPrefix("/user").Subrouter()
	api.RegisterHandlers(userRouter, http.MethodGet, "", nil, res.get)
	api.RegisterHandlers(userRouter, http.MethodGet, "/{id}", nil, res.getById)
	api.RegisterHandlers(userRouter, http.MethodPost, "", nil, res.add)

}

func (res *resource) get(w http.ResponseWriter, r *http.Request) {
	users, err := res.service.Get(r.Context())
	if err != nil {
		res.log.Error(err.Error())
		api.Write(w, http.StatusInternalServerError, api.NewResponse(false, "", err))
		return
	}
	res.log.Info("get all users")
	api.Write(w, http.StatusOK, users)
}

func (res *resource) getById(w http.ResponseWriter, r *http.Request) {
	user, err := res.service.GetUserByID(r.Context(), mux.Vars(r)["id"])
	if err != nil {
		api.Write(w, http.StatusInternalServerError, api.NewResponse(false, "", err))
		return
	}
	api.Write(w, http.StatusOK, user)
}

func (res *resource) add(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := api.SanitizeRequest(r, &user)
	if err != nil {
		api.Write(w, http.StatusInternalServerError, api.NewResponse(false, api.ErrorInvalidRequest, err))
		return

	}
	user.ID = uuid.New()
	err = res.service.Add(r.Context(), &user)
	if err != nil {
		api.Write(w, http.StatusInternalServerError, api.NewResponse(false, "failed to add user", err))
		return
	}
	api.Write(w, http.StatusOK, api.NewResponse(true, "user has been added", nil))
}
