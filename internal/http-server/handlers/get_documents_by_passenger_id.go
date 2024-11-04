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

func GetDocumentsByPassengerId(ctx context.Context, storage storage.Storage, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetDocumentsByPassengerId"

		w.Header().Set("Content-Type", "application/json")
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		passengerId := chi.URLParam(r, "passengerId")
		documents, err := storage.GetDocumentsByPassengerId(ctx, passengerId)
		if err != nil {
			log.Error("Error get documents by passenger id: ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.JSON(w, r, documents)
	}
}
