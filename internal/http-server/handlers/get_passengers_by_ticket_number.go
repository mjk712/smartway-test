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

// GetPassengersByTicketNumberHandler получает список пассажиров по номеру билета
// @Summary Получение списка пассажиров по номеру билета
// @Description Возвращает список пассажиров по номеру билета
// @Tags Пассажиры
// @Produce json
// @Param ticketNumber path string true "Номер билета"
// @Success 200 {array} models.Passenger "Список пассажиров успешно получен"
// @Failure 400 "Ошибка запроса или получения списка пассажиров"
// @Router /api/passengers/{ticketNumber} [get]
func GetPassengersByTicketNumberHandler(ctx context.Context, flightService service.FlightService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetPassengersByTicketNumber"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		ticketNumber := chi.URLParam(r, "ticketNumber")
		passengers, err := flightService.GetPassengersByTicketNumber(ctx, ticketNumber)
		if err != nil {
			log.Error("Error get passengers by ticket number: ", tools.ErrAttr(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, passengers)
	}
}
