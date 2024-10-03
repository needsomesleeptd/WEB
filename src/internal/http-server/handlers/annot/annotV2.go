package annot_handler

import (
	service "annotater/internal/bl/annotationService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

const (
	ClassTypeFieldNameV2 = "class_id"
	BbsFieldNameV2       = "bbs"
)

type RequestAddAnnotV2 struct { //other data is a file
	ErrorBB    []float32 `json:"error_bb"`
	ClassLabel uint64    `json:"class_label"`
}

type MarkUpsMetaData struct {
	ID         uint64    `json:"id"`
	ErrorBB    []float32 `json:"error_bb"`
	ClassLabel uint64    `json:"class_label"`
}

type ResponseGetAnnotsV2 struct {
	response.Response
	Markups []MarkUpsMetaData
}

func fromMarkupsSliceGetMeta(markups []models.Markup) []MarkUpsMetaData {
	metaDataSlice := make([]MarkUpsMetaData, len(markups))
	for i, markup := range markups {
		metaDataSlice[i] = MarkUpsMetaData{
			ID:         markup.ID,
			ErrorBB:    markup.ErrorBB,
			ClassLabel: markup.ClassLabel,
		}
	}
	return metaDataSlice
}

func NewAnnotHandlerV2(logSrc *logrus.Logger, servSrc service.IAnotattionService) AnnotHandlerV2 {
	return AnnotHandlerV2{
		log:          logSrc,
		annotService: servSrc,
	}
}

type AnnotHandlerV2 struct {
	log          *logrus.Logger
	annotService service.IAnotattionService
}

func (h *AnnotHandlerV2) AddAnnot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAddAnnot
		var pageData []byte
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		file, _, err := r.FormFile(AnnotFileFieldName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ErrorV2("invalid file")) //TODO:: add logging here
			h.log.Warn(err)
			return
		}
		pageData, err = io.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ErrorV2(models.GetUserError(err).Error())) //TODO:: add logging here
			h.log.Warn(err)
			return
		}
		classTypeString := r.FormValue(ClassTypeFieldNameV2)
		bbsString := r.FormValue(BbsFieldNameV2)

		err = json.Unmarshal([]byte(classTypeString), &req.ClassLabel)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ErrorV2(models.ErrDecodingRequest.Error())) //TODO:: add logging here
			h.log.Warn(err)
			return
		}
		fmt.Printf("values of bbs %v", bbsString)
		err = json.Unmarshal([]byte(bbsString), &req.ErrorBB)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ErrorV2(models.ErrDecodingRequest.Error())) //TODO:: add logging here
			h.log.Warn(err)
			return
		}

		annot := models.Markup{
			PageData:   pageData,
			ErrorBB:    req.ErrorBB,
			ClassLabel: req.ClassLabel,
			CreatorID:  userID,
		}
		err = h.annotService.AddAnottation(&annot)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//render.JSON(w, r, response.ErrorV2(models.GetUserError(err).Error()))
			h.log.Error(err)
			return
		}
		h.log.Infof("annot with class_label %v and bbs %v was successfully added", req.ClassLabel, req.ErrorBB)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *AnnotHandlerV2) GetAnnot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		annotIDStr := chi.URLParam(r, "id")
		annotUint64ID, err := strconv.ParseUint(annotIDStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ErrorV2(ErrDecodingRequest.Error()))
		}
		// Retrieve the annotation from the database or storage
		annot, err := h.annotService.GetAnottationByID(annotUint64ID)
		if err != nil {
			if !errors.Is(err, models.ErrNotFound) {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
			return
		}

		// Create a multipart writer
		markup := *models_dto.ToDtoMarkup(*annot)

		// Set response header for multipart/form-data
		w.Header().Set("Content-Type", "multipart/form-data")

		// Create a multipart writer
		writer := multipart.NewWriter(w)
		defer writer.Close()

		jsonData, err := json.Marshal(markup.ErrorBB)
		if err != nil {
			h.log.Warnf("annotV2 error : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = writer.WriteField(BbsFieldNameV2, string(jsonData))
		if err != nil {
			h.log.Warnf("annotV2 error : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		// Add class label
		strClassLabel := strconv.FormatUint(markup.ClassLabel, 10)
		err = writer.WriteField(ClassTypeFieldNameV2, strClassLabel)

		if err != nil {
			h.log.Warnf("annotV2 error : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		part, err := writer.CreateFormFile(AnnotFileFieldName, "file")
		if err != nil {
			h.log.Warnf("annotV2 error : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, err = part.Write(markup.PageData)
		if err != nil {
			h.log.Warnf("annotV2 error : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h.log.Infof("Annot with ID %v was successfully fetched", annotIDStr)

		w.WriteHeader(http.StatusOK)
		err = writer.Close()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.log.Warnf("annotV2 error : %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *AnnotHandlerV2) GetAllAnnots() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ErrorV2("Invalid JWT struct")) //TODO:: add logging here
			h.log.Warn("cannot get userIDfrom jwt in middleware")
			return
		}

		markUps, err := h.annotService.GetAllAnottations()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//render.JSON(w, r, response.ErrorV2(models.GetUserError(err).Error()))
			h.log.Warn(err)
			return
		}
		resp := ResponseGetAnnotsV2{Markups: fromMarkupsSliceGetMeta(markUps), Response: response.OK()}
		h.log.Infof("user with userID %v successfully got all annots\n", userID)
		render.JSON(w, r, resp)
	}
}

func (h *AnnotHandlerV2) DeleteAnnot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		annotIDStr := chi.URLParam(r, "id")
		annotUint64ID, err := strconv.ParseUint(annotIDStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.ErrorV2("invalid annotID"))
		}

		err = h.annotService.DeleteAnotattion(annotUint64ID)
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//render.JSON(w, r, response.ErrorV2(models.GetUserError(err).Error()))
			h.log.Warn(err)
			return
		}
		//h.log.Infof("user with userID %v successfully deleted annot\n", userID)
		w.WriteHeader(http.StatusOK)
	}
}
