package document_handler

import (
	service "annotater/internal/bl/documentService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	pdf_utils "annotater/internal/pkg/pdfUtils"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	CustomHeaderFilename = "filename"
)

type IDocumentHandlerV2 interface {
	LoadDocument(documentService service.IDocumentService) http.HandlerFunc
	GetDocumentsMetaData(documentService service.IDocumentService) http.HandlerFunc
	GetDocument(documentService service.IDocumentService) http.HandlerFunc
}

type DocumenthandlerV2 struct {
	logger     *logrus.Logger
	docService service.IDocumentService
}

func NewDocumentHandlerV2(logSrc *logrus.Logger, serv service.IDocumentService) DocumenthandlerV2 {
	return DocumenthandlerV2{
		logger:     logSrc,
		docService: serv,
	}
}

func (h *DocumenthandlerV2) GetDocumentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		documentID := chi.URLParam(r, "id")
		documentUUID, err := uuid.Parse(documentID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error()))
			return
		}
		document, err := h.docService.GetDocumentByID(documentUUID)

		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			h.logger.Error(err.Error())
			return
		}
		err = writeBytesIntoResponse(w, document.DocumentBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *DocumenthandlerV2) GetReportByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		documentID := chi.URLParam(r, "id")
		documentUUID, err := uuid.Parse(documentID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error()))
			return
		}
		report, err := h.docService.GetReportByID(documentUUID)

		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			h.logger.Error(err.Error())
			return
		}
		err = writeBytesIntoResponse(w, report.ReportData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			h.logger.Error(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *DocumenthandlerV2) GetDocumentsMetaData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		documentsMetaData, err := h.docService.GetDocumentsByCreatorID(userID)
		if err != nil {
			if errors.Is(err, models.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			h.logger.Error(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		resp := ResponseGettingMetaData{Response: response.OK(), DocumentsMetaData: documentsMetaData}
		render.JSON(w, r, resp)
	}
}

func (h *DocumenthandlerV2) CreateReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		pdfBytes, err := io.ReadAll(r.Body)
		defer r.Body.Close()

		var pagesCount int
		pagesCount, err = pdf_utils.GetPdfPageCount(pdfBytes)

		filename := r.Header.Get(CustomHeaderFilename)
		if err != nil {
			h.logger.Error(errors.Join(err, ErrGettingPageCount).Error())
			pagesCount = -1
		}

		documentMetaData := models.DocumentMetaData{
			ID:           uuid.New(),
			CreatorID:    userID,
			DocumentName: filename,
			CreationTime: time.Now(),
			PageCount:    pagesCount,
		}
		documentData := models.DocumentData{
			DocumentBytes: pdfBytes,
			ID:            documentMetaData.ID,
		}

		_, err = h.docService.LoadDocument(documentMetaData, documentData)
		if err != nil {
			if errors.Is(err, service.ErrDocumentFormat) {
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			h.logger.Error(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
