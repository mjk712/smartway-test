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

func DeleteTicketHandler(ctx context.Context, storage storage.Storage, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.DeleteTicket"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		ticketId := chi.URLParam(r, "ticketId")
		err := storage.DeleteTicketById(ctx, ticketId)
		if err != nil {
			log.Error("Error delete ticket with id: ", ticketId, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("ticket with id %s deleted", ticketId)))
	}
}