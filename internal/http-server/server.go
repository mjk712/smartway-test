package http_server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net"
	"net/http"
	"smartway-test/internal/config"
	"smartway-test/internal/http-server/handlers"
	"smartway-test/internal/storage"
	"time"
)

func NewServer(ctx context.Context, log *slog.Logger, cfg *config.Config, repo storage.Storage) *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api", func(r chi.Router) {
		r.Get("/tickets", handlers.GetTicketsHandler(ctx, repo, log))
		r.Get("/passengers/{ticketNumber}", handlers.GetPassengersByTicketNumberHandler(ctx, repo, log))
		r.Get("/documents/{passengerId}", handlers.GetDocumentsByPassengerId(ctx, repo, log))
		r.Get("/ticket/{ticketNumber}", handlers.GetTicketFullInfo(ctx, repo, log))
		r.Get("/reports/passenger/{passengerId}", handlers.GetPassengerReport(ctx, repo, log))
		r.Put("/ticket/{ticketId}", handlers.UpdateTicketInfo(ctx, repo, log))
		r.Put("/passenger/{passengerId}", handlers.UpdatePassengerInfo(ctx, repo, log))
		r.Put("/document/{documentId}", handlers.UpdateDocumentInfo(ctx, repo, log))
		r.Delete("/ticket/{ticketId}", handlers.DeleteTicketHandler(ctx, repo, log))
		r.Delete("/passenger/{passengerId}", handlers.DeletePassengerHandler(ctx, repo, log))
		r.Delete("/document/{documentId}", handlers.DeleteDocumentHandler(ctx, repo, log))

	})

	return &http.Server{
		Addr: cfg.Address,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		Handler: router,
	}
}
