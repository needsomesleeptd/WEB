package user_handler

import (
	service "annotater/internal/bl/userService"
	response "annotater/internal/lib/api"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type UserHandlerV2 struct {
	logger      *logrus.Logger
	userService service.IUserService
}

func NewUserHandlerV2(logSrc *logrus.Logger, serv service.IUserService) UserHandlerV2 {
	return UserHandlerV2{
		logger:      logSrc,
		userService: serv,
	}
}

type RequestChangeRoleV2 struct {
	ReqRole models.Role `json:"req_role"`
}

func (h *UserHandlerV2) ChangeUserPerms() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userLoginStr := chi.URLParam(r, "login")

		var req RequestChangeRoleV2
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error(ErrDecodingJson.Error()))
			h.logger.Warn(err.Error())
			return
		}

		err = h.userService.ChangeUserRoleByLogin(userLoginStr, req.ReqRole)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) || errors.Is(err, models.ErrInvalidRole) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			h.logger.Warn(err.Error())
			return
		}
		h.logger.Infof("successfully changed role of user with login %v  to role %v\n", userLoginStr, req.ReqRole)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *UserHandlerV2) GetAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		users, err := h.userService.GetAllUsers()
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Warn(err.Error())
			return
		}
		usersDTO := models_dto.ToDtoUserSlice(users)
		resp := ResponseGetAllUsers{response.OK(), usersDTO}
		h.logger.Infof("succesfully got all users\n")
		render.JSON(w, r, resp)
	}
}
