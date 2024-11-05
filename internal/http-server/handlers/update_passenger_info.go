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

// UpdatePassengerInfo обновляет информацию о пассажире
// @Summary Обновление информации о пассажире
// @Description Обновляет данные пассажира по заданному id
// @Tags Пассажиры
// @Accept json
// @Produce json
// @Param passengerId path string true "ID пассажира"
// @Param passenger body requests.UpdatePassengerRequest true "Данные для обновления информации о пассажире"
// @Success 200 {object} models.Passenger "Пассажир успешно обновлён"
// @Failure 400 "Ошибка запроса или обновления пассажира"
// @Router /api/passenger/{passengerId} [put]
func UpdatePassengerInfo(ctx context.Context, flightService service.FlightService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.UpdatePassengerInfo"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		passengerId := chi.URLParam(r, "passengerId")

		var req requests.UpdatePassengerRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("Error read request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updatedPassenger, err := flightService.UpdatePassengerInfo(ctx, passengerId, req)
		if err != nil {
			log.Error("Error updating ticket", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, updatedPassenger)
	}
}
