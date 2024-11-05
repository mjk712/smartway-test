package http_server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net"
	"net/http"
	_ "smartway-test/docs"
	"smartway-test/internal/config"
	"smartway-test/internal/http-server/handlers"
	"smartway-test/internal/service"
	"time"
)

func NewServer(ctx context.Context, log *slog.Logger, cfg *config.Config, flightService service.FlightService) *http.Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/api", func(r chi.Router) {
		r.Get("/tickets", handlers.GetTicketsHandler(ctx, flightService, log))
		r.Get("/passengers/{ticketNumber}", handlers.GetPassengersByTicketNumberHandler(ctx, flightService, log))
		r.Get("/documents/{passengerId}", handlers.GetDocumentsByPassengerId(ctx, flightService, log))
		r.Get("/ticket/{ticketNumber}", handlers.GetTicketFullInfo(ctx, flightService, log))
		r.Get("/reports/passenger/{passengerId}", handlers.GetPassengerReport(ctx, flightService, log))
		r.Put("/ticket/{ticketId}", handlers.UpdateTicketInfo(ctx, flightService, log))
		r.Put("/passenger/{passengerId}", handlers.UpdatePassengerInfo(ctx, flightService, log))
		r.Put("/document/{documentId}", handlers.UpdateDocumentInfo(ctx, flightService, log))
		r.Delete("/ticket/{ticketId}", handlers.DeleteTicketHandler(ctx, flightService, log))
		r.Delete("/passenger/{passengerId}", handlers.DeletePassengerHandler(ctx, flightService, log))
		r.Delete("/document/{documentId}", handlers.DeleteDocumentHandler(ctx, flightService, log))

		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("/api/swagger/doc.json"),
		))
	})

	return &http.Server{
		Addr: cfg.Address,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		Handler: router,
	}
}
