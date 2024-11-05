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

// UpdateTicketInfo обновляет информацию о билете
// @Summary Обновление информации о билете
// @Description Обновляет данные билета по заданному `ticketId`
// @Tags Билеты
// @Accept json
// @Produce json
// @Param ticketId path string true "ID билета"
// @Param ticket body requests.TicketUpdateRequest true "Данные для обновления информации о билете"
// @Success 200 {object} models.Ticket "Билет успешно обновлён"
// @Failure 400 "Ошибка запроса или обновления билета"
// @Router /ticket/{ticketId} [put]
func UpdateTicketInfo(ctx context.Context, flightService service.FlightService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.UpdateTicketInfo"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ticketId := chi.URLParam(r, "ticketId")

		var req requests.TicketUpdateRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("Error read request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedTicket, err := flightService.UpdateTicketInfo(ctx, ticketId, req)
		if err != nil {
			log.Error("Error updating ticket", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, updatedTicket)
	}
}
