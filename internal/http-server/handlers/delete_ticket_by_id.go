package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"smartway-test/internal/service"
)

// DeleteTicketHandler удаляет билет по id.
//
// @Summary Удаление билета
// @Description Удаляет билет по id и связь билета с пассажиром.
// @Tags Билеты
// @Accept json
// @Param ticketId path int true "Ticket ID"
// @Success 200 "Билет успешно удалён"
// @Failure 404 "Ошибка в запросе или при удалении билета"
// @Router /api/ticket/{ticketId} [delete]
func DeleteTicketHandler(ctx context.Context, flightService service.FlightService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.DeleteTicket"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ticketId := chi.URLParam(r, "ticketId")
		err := flightService.DeleteTicketById(ctx, ticketId)
		if err != nil {
			log.Error("Error delete ticket with id: ", ticketId, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, fmt.Sprintf("ticket with id %s deleted", ticketId))
	}
}
