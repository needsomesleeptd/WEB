package annot_type_handler

import (
	service "annotater/internal/bl/anotattionTypeService"
	response "annotater/internal/lib/api"
	logger_setup "annotater/internal/logger"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	auth_utils "annotater/internal/pkg/authUtils"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type RequestAnnotTypeV2 struct {
	Description string `json:"description"`
	ClassName   string `json:"class_name"`
}

type ResponseGetTypesV2 struct {
	MarkupTypes []models_dto.MarkupType
}

func NewAnnotTypehandlerV2(logSrc *logrus.Logger, servSrc service.IAnotattionTypeService) AnnotTypeHandlerV2 {
	return AnnotTypeHandlerV2{
		log:           logSrc,
		annotTypeServ: servSrc,
	}
}

type AnnotTypeHandlerV2 struct {
	annotTypeServ service.IAnotattionTypeService
	log           *logrus.Logger
}

func (h *AnnotTypeHandlerV2) AddAnnotType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAnnotTypeV2
		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			h.log.Warnf("cannot get userID from jwt %v in middleware", auth_utils.ExtractTokenFromReq(r))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("Invalud ID type")) //TODO:: add logging here
			return
		}
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.log.Warnf("unable to decode request with userID %v:%v", userID, err)
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}
		markupType := models.MarkupType{
			CreatorID:   int(userID),
			Description: req.Description,
			ClassName:   req.ClassName,
		}
		err = h.annotTypeServ.AddAnottationType(&markupType)
		if err != nil {
			if errors.Is(err, models.ErrDuplicateMarkupType) ||
				errors.Is(err, service.ErrInsertingEmptyClass) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			}
			h.log.Warn(err.Error())
			return
		}
		h.log.Infof("user with userID %v successfully added annotType with cls %v\n", userID, req.ClassName)
		//render.JSON(w, r, response.OK())
		w.WriteHeader(http.StatusOK)
	}
}

func (h *AnnotTypeHandlerV2) GetAllAnnotTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		markUpTypes, err := h.annotTypeServ.GetAllAnottationTypes()
		if err != nil {
			h.log.Warnf("unable to get all annot types %v\n", err.Error())
			//render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp := ResponseGetTypesV2{
			MarkupTypes: models_dto.ToDtoMarkupTypeSlice(markUpTypes),
		}
		h.log.Infof("successfully got annot types %v\n", resp.MarkupTypes)
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, resp)
	}
}

func (h *AnnotTypeHandlerV2) DeleteAnnotType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		annotTypeIDStr := chi.URLParam(r, "id")
		annotTypeUint64ID, err := strconv.ParseUint(annotTypeIDStr, 10, 64)
		if err != nil {
			h.log.Warnf(logger_setup.UnableToDecodeUserReqF, err.Error())
			render.JSON(w, r, response.Error(models.ErrDecodingRequest.Error()))
			w.WriteHeader(http.StatusBadRequest)
		}
		err = h.annotTypeServ.DeleteAnotattionType(annotTypeUint64ID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			h.log.Warnf("unable to delete annot type %v:%v\n", annotTypeUint64ID, err.Error())
			return
		}
		h.log.Infof("successfully deleted annot type %v\n", annotTypeUint64ID)
		w.WriteHeader(http.StatusOK)
	}
}
