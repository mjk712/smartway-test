package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"smartway-test/internal/storage"
)

func GetTicketFullInfo(ctx context.Context, storage storage.Storage, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetTicketFullInfo"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		ticketNumber := chi.URLParam(r, "ticketNumber")
		ticket, err := storage.GetFullTicketInfo(ctx, ticketNumber)
		if err != nil {
			log.Error("Error get passengers by ticket number: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, ticket)
	}
}
