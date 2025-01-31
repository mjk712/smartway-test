package storage

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"smartway-test/internal/http-server/requests"
	"smartway-test/internal/http-server/response"
	"smartway-test/internal/models"
	"smartway-test/internal/storage/query"
	"time"
)

type Storage interface {
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

type StorageRepo struct {
	DB *sqlx.DB
}

func New(connectionString string) (Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	//path for docker file:///app/internal/storage/migrations
	m, err := migrate.New("file:///app/internal/storage/migrations", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &StorageRepo{DB: db}, nil
}

func (s *StorageRepo) GetTickets(ctx context.Context) ([]models.Ticket, error) {
	const op = "storage.postgresql.GetTickets"

	var tickets []models.Ticket

	rows, err := s.DB.QueryxContext(ctx, query.GetTickets)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var t models.Ticket
		if err := rows.StructScan(&t); err != nil {
			return nil, fmt.Errorf("scan rows %s: %w", op, err)
		}
		tickets = append(tickets, t)
	}
	rows.Close()
	return tickets, nil
}

func (s *StorageRepo) GetPassengersByTicketNumber(ctx context.Context, ticketNumber string) ([]models.Passenger, error) {
	const op = "storage.postgresql.GetPassengersByTicketNumber"
	var passengers []models.Passenger

	rows, err := s.DB.QueryxContext(ctx, query.GetPassengersByTicketNumber, ticketNumber)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}
	for rows.Next() {
		var p models.Passenger
		if err := rows.StructScan(&p); err != nil {
			return nil, fmt.Errorf("scan rows %s: %w", op, err)
		}
		passengers = append(passengers, p)
	}
	rows.Close()
	return passengers, nil
}

func (s *StorageRepo) GetDocumentsByPassengerId(ctx context.Context, passengerId string) ([]models.Document, error) {
	const op = "storage.postgresql.GetDocumentsByPassengerId"
	var documents []models.Document
	rows, err := s.DB.QueryxContext(ctx, query.GetDocumentsByPassengerId, passengerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var d models.Document
		if err := rows.StructScan(&d); err != nil {
			return nil, fmt.Errorf("scan rows %s: %w", op, err)
		}
		documents = append(documents, d)
	}
	rows.Close()
	return documents, nil
}

func (s *StorageRepo) GetFullTicketInfo(ctx context.Context, ticketNumber string) (response.FullTicketInfo, error) {
	const op = "storage.postgresql.GetFullTicketInfo"

	type combinedRow struct {
		models.Ticket
		models.Passenger
		models.Document
	}

	var combinedRows []combinedRow

	if err := s.DB.SelectContext(ctx, &combinedRows, query.GetTicketFullInfo, ticketNumber); err != nil {
		return response.FullTicketInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	var fullTicket response.FullTicketInfo
	ticketExist := false
	passengerMap := make(map[int]map[int]*response.PassengerWithDocs)

	for _, row := range combinedRows {

		if !ticketExist {
			fullTicket = response.FullTicketInfo{
				Ticket:     row.Ticket,
				Passengers: []response.PassengerWithDocs{},
			}
			passengerMap[row.TicketId] = make(map[int]*response.PassengerWithDocs)
			ticketExist = true
		}
		//проверяем наличие пассажира
		passengerInfo, passengerExists := passengerMap[row.TicketId][row.Passenger.PassengerId]
		if !passengerExists {
			//добавляем пассажира
			passengerMap[row.TicketId][row.Passenger.PassengerId] = &response.PassengerWithDocs{
				Passenger: row.Passenger,
				Documents: []models.Document{},
			}
			passengerInfo = passengerMap[row.TicketId][row.Passenger.PassengerId]
			fullTicket.Passengers = append(fullTicket.Passengers, *passengerInfo)
		}

		//проверяем документ
		if row.DocumentId != 0 && row.Document.PassengerId == row.Passenger.PassengerId {

			for i := range fullTicket.Passengers {
				if fullTicket.Passengers[i].PassengerId == row.Passenger.PassengerId {
					fullTicket.Passengers[i].Documents =
						append(fullTicket.Passengers[i].Documents, row.Document)
				}
			}
		}

	}

	return fullTicket, nil
}

func (s *StorageRepo) UpdateTicketInfo(ctx context.Context, ticketId string, updateData requests.TicketUpdateRequest) (models.Ticket, error) {
	const op = "storage.postgresql.UpdateTicketInfo"

	queryBuilder := sq.Update("flight_ticket").
		Where(sq.Eq{"ticket_id": ticketId}).
		Suffix("RETURNING ticket_id, departure_point, destination_point, order_number, service_provider, departure_date, arrival_date, passenger_id AS \"ticket.passenger_id\", created_at").
		PlaceholderFormat(sq.Dollar)

	if updateData.DeparturePoint != nil {
		queryBuilder = queryBuilder.Set("departure_point", updateData.DeparturePoint)
	}
	if updateData.DestinationPoint != nil {
		queryBuilder = queryBuilder.Set("destination_point", updateData.DestinationPoint)
	}
	if updateData.ServiceProvider != nil {
		queryBuilder = queryBuilder.Set("service_provider", updateData.ServiceProvider)
	}
	if updateData.DepartureDate != nil {
		queryBuilder = queryBuilder.Set("departure_date", updateData.DepartureDate)
	}
	if updateData.ArrivalDate != nil {
		queryBuilder = queryBuilder.Set("arrival_date", updateData.ArrivalDate)
	}
	if updateData.PassengerId != nil {
		queryBuilder = queryBuilder.Set("passenger_id", updateData.PassengerId)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return models.Ticket{}, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var updatedTicket models.Ticket
	if err := s.DB.QueryRowxContext(ctx, query, args...).StructScan(&updatedTicket); err != nil {
		return models.Ticket{}, fmt.Errorf("%s: %w", op, err)
	}
	return updatedTicket, nil
}

func (s *StorageRepo) UpdatePassengerInfo(ctx context.Context, passengerId string, updateData requests.UpdatePassengerRequest) (models.Passenger, error) {
	const op = "storage.postgresql.UpdatePassengerInfo"

	queryBuilder := sq.Update("passenger").
		Where(sq.Eq{"passenger_id": passengerId}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)

	if updateData.MiddleName != nil {
		queryBuilder = queryBuilder.Set("middle_name", *updateData.MiddleName)
	}
	if updateData.LastName != nil {
		queryBuilder = queryBuilder.Set("last_name", *updateData.LastName)
	}
	if updateData.FirstName != nil {
		queryBuilder = queryBuilder.Set("first_name", *updateData.FirstName)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return models.Passenger{}, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var updatedPassenger models.Passenger

	if err := s.DB.QueryRowxContext(ctx, query, args...).StructScan(&updatedPassenger); err != nil {
		return models.Passenger{}, fmt.Errorf("%s: %w", op, err)
	}

	return updatedPassenger, nil
}

func (s *StorageRepo) UpdateDocumentInfo(ctx context.Context, documentId string, updateData requests.DocumentUpdateRequest) (models.Document, error) {
	const op = "storage.postgresql.UpdateDocumentInfo"
	queryBuilder := sq.Update("document").
		Where(sq.Eq{"document_id": documentId}).
		Suffix("RETURNING document_id, passenger_id AS \"document.passenger_id\", document_type, document_number").
		PlaceholderFormat(sq.Dollar)

	if updateData.DocumentType != nil {
		queryBuilder = queryBuilder.Set("document_type", *updateData.DocumentType)
	}

	if updateData.DocumentNumber != nil {
		queryBuilder = queryBuilder.Set("document_number", *updateData.DocumentNumber)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return models.Document{}, fmt.Errorf("%s: failed to build query: %w", op, err)
	}

	var updatedDocument models.Document
	if err := s.DB.QueryRowxContext(ctx, query, args...).StructScan(&updatedDocument); err != nil {
		return models.Document{}, fmt.Errorf("%s: %w", op, err)
	}
	return updatedDocument, nil
}

func (s *StorageRepo) DeleteTicketById(ctx context.Context, ticketId string) error {
	const op = "storage.postgres.DeleteTicketById"

	_, err := s.DB.QueryxContext(ctx, query.DeleteTicket, ticketId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *StorageRepo) DeletePassengerById(ctx context.Context, passengerId string) error {
	const op = "storage.postgres.DeletePassengerById"

	_, err := s.DB.QueryxContext(ctx, query.DeletePassenger, passengerId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *StorageRepo) DeleteDocumentById(ctx context.Context, documentId string) error {
	const op = "storage.postgres.DeleteDocumentById"

	_, err := s.DB.QueryxContext(ctx, query.DeleteDocument, documentId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *StorageRepo) GetPassengerReport(ctx context.Context, passengerId string, startDate time.Time, endDate time.Time) ([]response.FlightReport, error) {
	const op = "storage.postgres.GetPassengerReport"

	var report []response.FlightReport
	err := s.DB.SelectContext(ctx, &report, query.GetPassengerReport, startDate, endDate, passengerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return report, nil
}
