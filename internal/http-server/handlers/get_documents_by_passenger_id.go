package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"smartway-test/internal/service"
	"smartway-test/internal/tools"
)

// GetDocumentsByPassengerId Возвращает список документов по id пассажира.
//
// @Summary Получение списка документов
// @Description Возвращает список документов. Требует id пассажира.
// @Tags Документы
// @Produce json
// @Param passengerId path int true "Passenger ID"
// @Success 200 {array} models.Document "Список продуктов успешно получен"
// @Failure 404 "Ошибка в запросе или при получении списка документов"
// @Router /api/documents/{passengerId} [get]
func GetDocumentsByPassengerId(ctx context.Context, flightService service.FlightService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetDocumentsByPassengerId"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		passengerId := chi.URLParam(r, "passengerId")
		documents, err := flightService.GetDocumentsByPassengerId(ctx, passengerId)
		if err != nil {
			log.Error("Error get documents by passenger id: ", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, documents)
	}
}
