package handlers

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"smartway-test/internal/storage"
)

func DeletePassengerHandler(ctx context.Context, storage storage.Storage, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.DeletePassengerHandler"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		passengerId := chi.URLParam(r, "passengerId")
		err := storage.DeletePassengerById(ctx, passengerId)
		if err != nil {
			log.Error("Error delete ticket with id: ", passengerId, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("passenger with id %s deleted", passengerId)))
	}
}
