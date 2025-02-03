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

// GetTicketFullInfo получает полную информацию о билете
// @Summary Получение полной информации о билете
// @Description Возвращает полные данные о билете по заданному номеру
// @Tags Билеты
// @Produce json
// @Param ticketNumber path string true "Номер билета"
// @Success 200 {object} response.FullTicketInfo "Информация о билете успешно получена"
// @Failure 400 "Ошибка запроса или получения полной информации о билете"
// @Router /api/ticket/{ticketNumber} [get]
func GetTicketFullInfo(ctx context.Context, flightService service.FlightService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetTicketFullInfo"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		ticketNumber := chi.URLParam(r, "ticketNumber")
		ticket, err := flightService.GetFullTicketInfo(ctx, ticketNumber)
		if err != nil {
			log.Error("Error get passengers by ticket number: ", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, ticket)
	}
}
