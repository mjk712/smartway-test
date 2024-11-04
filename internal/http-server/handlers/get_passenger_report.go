package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"smartway-test/internal/storage"
	"time"
)

// GetPassengerReport получает отчет о пассажире за определенный период
// @Summary Получение отчета о пассажире
// @Description Возвращает отчет о пассажире по заданному `passengerId` и диапазону дат
// @Tags Пассажиры
// @Produce json
// @Param passengerId path int true "ID пассажира"
// @Param start_date query string true "Дата начала в формате YYYY-MM-DD"
// @Param end_date query string true "Дата окончания в формате YYYY-MM-DD"
// @Success 200 {array} response.FlightReport "Отчёт успешно получен"
// @Failure 400 "Ошибка запроса или получения отчёта"
// @Router /reports/passenger/{passengerId} [get]
func GetPassengerReport(ctx context.Context, storage storage.Storage, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		passengerId := chi.URLParam(r, "passengerId")
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")

		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			log.Error("Invalid start date", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			log.Error("Invalid end date", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		report, err := storage.GetPassengerReport(ctx, passengerId, startDate, endDate)
		if err != nil {
			log.Error("Error getting report", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		render.JSON(w, r, &report)
	}
}
