package auth_handler

import (
	auth_service "annotater/internal/bl/auth"
	response "annotater/internal/lib/api"
	logger_setup "annotater/internal/logger"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

var (
	ErrDecodingJson = errors.New("invalid json")
	ErrInternalServ = errors.New("internal server error")
)

const (
	COOKIE_NAME = "auth_jwt"
)

type RequestSignUp struct {
	User models_dto.UserWithPasswd `json:"user"`
}

type RequestSignIn struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ResponseSignIn struct {
	Response response.Response
	Jwt      string `json:"jwt,omitempty"`
}

type AuthHandlerV1 struct {
	log         *logrus.Logger
	authService auth_service.IAuthService
}

func NewAuthHandlerV1(logSrc *logrus.Logger, authServSrc auth_service.IAuthService) AuthHandlerV1 {
	return AuthHandlerV1{
		log:         logSrc,
		authService: authServSrc,
	}
}

func (h *AuthHandlerV1) SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSignUp
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.log.Warnf(logger_setup.UnableToDecodeUserReqF, err)
			render.JSON(w, r, response.Error(ErrDecodingJson.Error())) //TODO:: add logging here
			return
		}
		req.User.Role = models.Sender
		userDto := models_dto.FromDtoUserWithPasswd(&req.User)
		candidate := models_dto.FromDtoUser(&userDto)
		h.log.Infof("candidate : %v, req: %v", candidate, req)
		err = h.authService.SignUp(&candidate)
		if err != nil {
			h.log.Warnf("unable to signUp with user login %v:%v\n", req.User.Login, err)
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}

		h.log.Infof("user with login %v successfuly signed up\n", req.User.Login)
		render.JSON(w, r, response.OK())
	}
}

func (h *AuthHandlerV1) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestSignIn
		var tokenStr string
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.log.Warnf(logger_setup.UnableToDecodeUserReqF, err)
			render.JSON(w, r, ResponseSignIn{Response: response.Error(ErrDecodingJson.Error())})
			return
		}
		candidate := models.User{Login: req.Login, Password: req.Password}
		tokenStr, err = h.authService.SignIn(&candidate)
		if err != nil {
			h.log.Warnf("unable to signIn with user login %v:%v\n", req.Login, err)
			render.JSON(w, r, ResponseSignIn{Response: response.Error(models.GetUserError(err).Error())})
			return
		}

		resp := ResponseSignIn{Response: response.OK(), Jwt: tokenStr}
		h.log.Infof("user with login %v sucessfully signed in\n", req.Login)
		render.JSON(w, r, resp)
	}
}
