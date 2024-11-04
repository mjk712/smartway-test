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
	"smartway-test/internal/storage"
)

func UpdateTicketInfo(ctx context.Context, storage storage.Storage, log *slog.Logger) http.HandlerFunc {
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

		updatedTicket, err := storage.UpdateTicketInfo(ctx, ticketId, req)
		if err != nil {
			log.Error("Error updating ticket", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, updatedTicket)
	}
}
