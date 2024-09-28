package auth_handler

import (
	auth_service "annotater/internal/bl/auth"
	service "annotater/internal/bl/auth"
	response "annotater/internal/lib/api"
	logger_setup "annotater/internal/logger"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type RequestSignUpV2 struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Group    string `json:"group"`
}

func fromSignIpRequestUser(req *RequestSignUpV2) models_dto.User {
	return models_dto.User{
		Login:    req.Login,
		Password: req.Password,
		Name:     req.Name,
		Surname:  req.Surname,
		Group:    req.Group,
	}
}

type AuthHandlerV2 struct {
	log         *logrus.Logger
	authService auth_service.IAuthService
}

func NewAuthHandlerV2(logSrc *logrus.Logger, authServSrc auth_service.IAuthService) AuthHandlerV2 {
	return AuthHandlerV2{
		log:         logSrc,
		authService: authServSrc,
	}
}

func (h *AuthHandlerV2) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSignUpV2
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.log.Warnf(logger_setup.UnableToDecodeUserReqF, err)
			render.JSON(w, r, response.Error(ErrDecodingJson.Error())) //TODO:: add logging here
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		candidateDto := fromSignIpRequestUser(&req)
		candidate := models_dto.FromDtoUser(&candidateDto)
		err = h.authService.SignUp(&candidate)
		if err != nil {
			h.log.Warnf("unable to signUp with user login %v:%v\n", candidate.Login, err)
			if errors.Is(err, models.ErrDuplicateuserData) ||
				errors.Is(err, auth_service.ErrNoLogin) ||
				errors.Is(err, auth_service.ErrNoPasswd) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		h.log.Infof("user with login %v successfuly signed up\n", candidate.Login)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *AuthHandlerV2) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSignIn
		var tokenStr string
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.log.Warnf(logger_setup.UnableToDecodeUserReqF, err)
			render.JSON(w, r, ResponseSignIn{Response: response.Error(ErrDecodingJson.Error())})
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		candidate := models.User{Login: req.Login, Password: req.Password}
		tokenStr, err = h.authService.SignIn(&candidate)

		if err != nil {
			h.log.Warnf("unable to signIn with user login %v:%v\n", req.Login, err)

			if errors.Is(err, auth_service.ErrNoLogin) ||
				errors.Is(err, auth_service.ErrNoPasswd) ||
				errors.Is(err, models.ErrNotFound) ||
				errors.Is(err, service.ErrWrongPassword) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		resp := ResponseSignIn{Response: response.OK(), Jwt: tokenStr}
		h.log.Infof("user with login %v sucessfully signed in\n", req.Login)
		render.JSON(w, r, resp)
		w.WriteHeader(http.StatusOK)
	}
}
