package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"smartway-test/internal/service"
)

// DeleteDocumentHandler удаляет документ по id.
//
// @Summary Удаление документа
// @Description Удаляет документ по id.
// @Tags Документы
// @Accept json
// @Param documentId path int true "Document ID"
// @Success 200 "Документ успешно удалён"
// @Failure 404 "Ошибка в запросе или при удалении документа"
// @Router /api/document/{documentId} [delete]
func DeleteDocumentHandler(ctx context.Context, flightService service.FlightService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.DeleteDocumentHandler"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		documentId := chi.URLParam(r, "documentId")
		err := flightService.DeleteDocumentById(ctx, documentId)
		if err != nil {
			log.Error("Error delete document with id: ", documentId, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("document with id %s deleted", documentId)))
	}
}
