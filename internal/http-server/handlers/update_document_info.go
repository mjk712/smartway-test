package handlers

import (
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"smartway-test/internal/http-server/requests"
	"smartway-test/internal/service"
)

// UpdateDocumentInfo обновляет информацию о документе
// @Summary Обновление информации о документе
// @Description Обновляет данные документа по id
// @Tags Документы
// @Accept json
// @Produce json
// @Param documentId path string true "ID документа"
// @Param document body requests.DocumentUpdateRequest true "Данные для обновления документа"
// @Success 200 {object} models.Document "Документ успешно обновлён"
// @Failure 400 "Ошибка запроса или обновления документа"
// @Router /document/{documentId} [put]
func UpdateDocumentInfo(ctx context.Context, flightService service.FlightService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.UpdateDocumentInfo"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		documentId := chi.URLParam(r, "documentId")

		var req requests.DocumentUpdateRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("Error read request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedDocument, err := flightService.UpdateDocumentInfo(ctx, documentId, req)
		if err != nil {
			log.Error("Error updating ticket", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, updatedDocument)
	}
}
