package service

import (
	"context"
	"fmt"
	"smartway-test/internal/http-server/requests"
	"smartway-test/internal/http-server/response"
	"smartway-test/internal/models"
	"smartway-test/internal/storage"
	"time"
)

type FlightService interface {
	GetTickets(ctx context.Context) ([]models.Ticket, error)
	GetPassengersByTicketNumber(ctx context.Context, ticketNumber string) ([]models.Passenger, error)
	GetDocumentsByPassengerId(ctx context.Context, passengerId string) ([]models.Document, error)
	GetFullTicketInfo(ctx context.Context, ticketNumber string) (response.FullTicketInfo, error)
	GetPassengerReport(ctx context.Context, passengerId string, startDate time.Time, endDate time.Time) ([]response.FlightReport, error)
	UpdateTicketInfo(ctx context.Context, ticketId string, updateData requests.TicketUpdateRequest) (models.Ticket, error)
	UpdatePassengerInfo(ctx context.Context, passengerId string, updateData requests.UpdatePassengerRequest) (models.Passenger, error)
	UpdateDocumentInfo(ctx context.Context, documentId string, updateData requests.DocumentUpdateRequest) (models.Document, error)
	DeleteTicketById(ctx context.Context, ticketId string) error
	DeletePassengerById(ctx context.Context, passengerId string) error
	DeleteDocumentById(ctx context.Context, documentId string) error
}

type FlightServiceImpl struct {
	storageRepository storage.Storage
}

func NewFlightService(repo storage.Storage) FlightService {
	return &FlightServiceImpl{
		storageRepository: repo,
	}
}

func (f *FlightServiceImpl) GetTickets(ctx context.Context) ([]models.Ticket, error) {
	const op = "FlightService.GetTickets"
	tickets, err := f.storageRepository.GetTickets(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tickets, nil
}

func (f *FlightServiceImpl) GetPassengersByTicketNumber(ctx context.Context, ticketNumber string) ([]models.Passenger, error) {
	const op = "FlightService.GetPassengersByTicketNumber"
	passengers, err := f.storageRepository.GetPassengersByTicketNumber(ctx, ticketNumber)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return passengers, nil
}

func (f *FlightServiceImpl) GetDocumentsByPassengerId(ctx context.Context, passengerId string) ([]models.Document, error) {
	const op = "FlightService.GetDocumentsByPassengerId"
	documents, err := f.storageRepository.GetDocumentsByPassengerId(ctx, passengerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return documents, nil
}

func (f *FlightServiceImpl) GetFullTicketInfo(ctx context.Context, ticketNumber string) (response.FullTicketInfo, error) {
	const op = "FlightService.GetFullTicketInfo"
	ticketInfo, err := f.storageRepository.GetFullTicketInfo(ctx, ticketNumber)
	if err != nil {
		return response.FullTicketInfo{}, fmt.Errorf("%s: %w", op, err)
	}
	return ticketInfo, nil
}

func (f *FlightServiceImpl) GetPassengerReport(ctx context.Context, passengerId string, startDate time.Time, endDate time.Time) ([]response.FlightReport, error) {
	const op = "FlightService.GetPassengerReport"
	report, err := f.storageRepository.GetPassengerReport(ctx, passengerId, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return report, nil
}

func (f *FlightServiceImpl) UpdateTicketInfo(ctx context.Context, ticketId string, updateData requests.TicketUpdateRequest) (models.Ticket, error) {
	const op = "FlightService.UpdateTicketInfo"
	updatedTicket, err := f.storageRepository.UpdateTicketInfo(ctx, ticketId, updateData)
	if err != nil {
		return updatedTicket, fmt.Errorf("%s: %w", op, err)
	}
	return updatedTicket, nil
}

func (f *FlightServiceImpl) UpdatePassengerInfo(ctx context.Context, passengerId string, updateData requests.UpdatePassengerRequest) (models.Passenger, error) {
	const op = "FlightService.UpdatePassengerInfo"
	passenger, err := f.storageRepository.UpdatePassengerInfo(ctx, passengerId, updateData)
	if err != nil {
		return passenger, fmt.Errorf("%s: %w", op, err)
	}
	return passenger, nil
}

func (f *FlightServiceImpl) UpdateDocumentInfo(ctx context.Context, documentId string, updateData requests.DocumentUpdateRequest) (models.Document, error) {
	const op = "FlightService.UpdateDocumentInfo"
	updatedDocument, err := f.storageRepository.UpdateDocumentInfo(ctx, documentId, updateData)
	if err != nil {
		return updatedDocument, fmt.Errorf("%s: %w", op, err)
	}
	return updatedDocument, nil
}

func (f *FlightServiceImpl) DeleteTicketById(ctx context.Context, ticketId string) error {
	const op = "FlightService.DeleteTicketById"
	err := f.storageRepository.DeleteTicketById(ctx, ticketId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (f *FlightServiceImpl) DeletePassengerById(ctx context.Context, passengerId string) error {
	const op = "FlightService.DeletePassengerById"
	err := f.storageRepository.DeletePassengerById(ctx, passengerId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (f *FlightServiceImpl) DeleteDocumentById(ctx context.Context, documentId string) error {
	const op = "FlightService.DeleteDocumentById"
	err := f.storageRepository.DeleteDocumentById(ctx, documentId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
