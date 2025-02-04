package handlers

import (
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"smartway-test/internal/service"
	"smartway-test/internal/tools"
)

// GetTicketsHandler получает список всех билетов
// @Summary Получение списка билетов
// @Description Возвращает список всех доступных билетов
// @Tags Билеты
// @Produce json
// @Success 200 {array} models.Ticket "Список билетов успешно получен"
// @Failure 400 "Ошибка запроса или получения списка билетов"
// @Router /api/tickets [get]
func GetTicketsHandler(ctx context.Context, flightService service.FlightService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetTickets"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		tickets, err := flightService.GetTickets(ctx)
		if err != nil {
			log.Error("Error get tickets: ", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, tickets)
	}
}
