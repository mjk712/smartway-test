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

// DeletePassengerHandler удаляет пассажира по id.
//
// @Summary Удаление пассажира
// @Description Удаляет пассажира по id.
// @Tags Пассажиры
// @Accept json
// @Param passengerId path int true "Passenger ID"
// @Success 200 "Пассажир успешно удалён"
// @Failure 404 "Ошибка в запросе или при удалении пассажира"
// @Router /api/passenger/{passengerId} [delete]
func DeletePassengerHandler(ctx context.Context, flightService service.FlightService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.DeletePassengerHandler"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		passengerId := chi.URLParam(r, "passengerId")
		err := flightService.DeletePassengerById(ctx, passengerId)
		if err != nil {
			log.Error("Error delete passenger with id: ", passengerId, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, fmt.Sprintf("passenger with id %s deleted", passengerId))
	}
}
