package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"smartway-test/internal/http-server/requests"
	"smartway-test/internal/http-server/response"
	"smartway-test/internal/models"
	"smartway-test/internal/storage/query"
	"strconv"
	"time"
)

type Storage interface {
	GetTickets(ctx context.Context) ([]*models.Ticket, error)
	GetPassengersByTicketNumber(ctx context.Context, ticketNumber string) ([]*models.Passenger, error)
	GetDocumentsByPassengerId(ctx context.Context, passengerId string) ([]*models.Document, error)
	GetFullTicketInfo(ctx context.Context, ticketNumber string) ([]*response.FullTicketInfo, error)
	GetPassengerReport(ctx context.Context, passengerId string, startDate time.Time, endDate time.Time) ([]*response.FlightReport, error)
	UpdateTicketInfo(ctx context.Context, ticketId string, updateData requests.TicketUpdateRequest) (*models.Ticket, error)
	UpdatePassengerInfo(ctx context.Context, passengerId string, updateData requests.UpdatePassengerRequest) (*models.Passenger, error)
	UpdateDocumentInfo(ctx context.Context, documentId string, updateData requests.DocumentUpdateRequest) (*models.Document, error)
	DeleteTicketById(ctx context.Context, ticketId string) error
	DeletePassengerById(ctx context.Context, passengerId string) error
	DeleteDocumentById(ctx context.Context, documentId string) error
}

type StorageRepo struct {
	db *sqlx.DB
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
	return &StorageRepo{db: db}, nil
}

func (s *StorageRepo) GetTickets(ctx context.Context) ([]*models.Ticket, error) {
	const op = "storage.postgresql.GetTickets"

	var tickets []*models.Ticket

	rows, err := s.db.QueryxContext(ctx, query.GetTickets)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var t models.Ticket
		if err := rows.StructScan(&t); err != nil {
			return nil, fmt.Errorf("scan rows %s: %w", op, err)
		}
		tickets = append(tickets, &t)
	}
	rows.Close()
	return tickets, nil
}

func (s *StorageRepo) GetPassengersByTicketNumber(ctx context.Context, ticketNumber string) ([]*models.Passenger, error) {
	const op = "storage.postgresql.GetPassengersByTicketNumber"
	var passengers []*models.Passenger

	rows, err := s.db.QueryxContext(ctx, query.GetPassengersByTicketNumber, ticketNumber)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}
	for rows.Next() {
		var p models.Passenger
		if err := rows.StructScan(&p); err != nil {
			return nil, fmt.Errorf("scan rows %s: %w", op, err)
		}
		passengers = append(passengers, &p)
	}
	rows.Close()
	return passengers, nil
}

func (s *StorageRepo) GetDocumentsByPassengerId(ctx context.Context, passengerId string) ([]*models.Document, error) {
	const op = "storage.postgresql.GetDocumentsByPassengerId"
	var documents []*models.Document
	rows, err := s.db.QueryxContext(ctx, query.GetDocumentsByPassengerId, passengerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	for rows.Next() {
		var d models.Document
		if err := rows.StructScan(&d); err != nil {
			return nil, fmt.Errorf("scan rows %s: %w", op, err)
		}
		documents = append(documents, &d)
	}
	rows.Close()
	return documents, nil
}

func (s *StorageRepo) GetFullTicketInfo(ctx context.Context, ticketNumber string) ([]*response.FullTicketInfo, error) {
	const op = "storage.postgresql.GetFullTicketInfo"

	type combinedRow struct {
		models.Ticket
		models.Passenger
		models.Document
	}

	var combinedRows []combinedRow

	if err := s.db.SelectContext(ctx, &combinedRows, query.GetTicketFullInfo, ticketNumber); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ticketMap := make(map[int]*response.FullTicketInfo)
	passengerMap := make(map[int]map[int]*response.PassengerWithDocs)

	for _, row := range combinedRows {
		ticketInfo, exists := ticketMap[row.TicketId]
		if !exists {
			//если билет новый, добавляем в мапу
			ticketMap[row.TicketId] = &response.FullTicketInfo{
				Ticket:     row.Ticket,
				Passengers: []response.PassengerWithDocs{},
			}
			ticketInfo = ticketMap[row.TicketId]
			passengerMap[row.TicketId] = make(map[int]*response.PassengerWithDocs)
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
			ticketInfo.Passengers = append(ticketInfo.Passengers, *passengerInfo)
		}

		//проверяем документ
		if row.DocumentId != 0 && row.Document.PassengerId == row.Passenger.PassengerId {

			for i := range ticketMap[row.TicketId].Passengers {
				if ticketMap[row.TicketId].Passengers[i].PassengerId == row.Passenger.PassengerId {
					ticketMap[row.TicketId].Passengers[i].Documents =
						append(ticketMap[row.TicketId].Passengers[i].Documents, row.Document)
				}
			}
		}

	}

	var result []*response.FullTicketInfo
	for _, info := range ticketMap {
		result = append(result, info)
	}

	return result, nil
}

func (s *StorageRepo) UpdateTicketInfo(ctx context.Context, ticketId string, updateData requests.TicketUpdateRequest) (*models.Ticket, error) {
	const op = "storage.postgresql.UpdateTicketInfo"

	query := "UPDATE flight_ticket SET"
	params := []interface{}{}
	idx := 1
	if updateData.DeparturePoint != nil {
		if idx > 1 {
			query += ", "
		}
		query += " departure_point = $" + strconv.Itoa(idx)
		params = append(params, *updateData.DeparturePoint)
		idx++
	}
	if updateData.DestinationPoint != nil {
		if idx > 1 {
			query += ", "
		}
		query += " destination_point = $" + strconv.Itoa(idx)
		params = append(params, *updateData.DestinationPoint)
		idx++
	}
	if updateData.ServiceProvider != nil {
		if idx > 1 {
			query += ", "
		}
		query += " service_provider = $" + strconv.Itoa(idx)
		params = append(params, *updateData.ServiceProvider)
		idx++
	}
	if updateData.DepartureDate != nil {
		if idx > 1 {
			query += ", "
		}
		query += " departure_date = $" + strconv.Itoa(idx)
		params = append(params, *updateData.DepartureDate)
		idx++
	}
	if updateData.ArrivalDate != nil {
		if idx > 1 {
			query += ", "
		}
		query += " arrival_date = $" + strconv.Itoa(idx)
		params = append(params, *updateData.ArrivalDate)
		idx++
	}
	query += " WHERE ticket_id = $" + strconv.Itoa(idx)
	idx++
	params = append(params, ticketId)
	query += " RETURNING *;"

	var updatedTicket models.Ticket
	if err := s.db.QueryRowxContext(ctx, query, params...).StructScan(&updatedTicket); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &updatedTicket, nil
}

func (s *StorageRepo) UpdatePassengerInfo(ctx context.Context, passengerId string, updateData requests.UpdatePassengerRequest) (*models.Passenger, error) {
	const op = "storage.postgresql.UpdatePassengerInfo"

	query := "UPDATE passenger SET"
	params := []interface{}{}
	idx := 1
	if updateData.FirstName != nil {
		if idx > 1 {
			query += ", "
		}
		query += " first_name = $" + strconv.Itoa(idx)
		params = append(params, *updateData.FirstName)
		idx++
	}
	if updateData.LastName != nil {
		if idx > 1 {
			query += ", "
		}
		query += " last_name = $" + strconv.Itoa(idx)
		params = append(params, *updateData.LastName)
		idx++
	}
	if updateData.MiddleName != nil {
		if idx > 1 {
			query += ", "
		}
		query += " middle_name = $" + strconv.Itoa(idx)
		params = append(params, *updateData.MiddleName)
		idx++
	}

	query += " WHERE passenger_id = $" + strconv.Itoa(idx)
	idx++
	params = append(params, passengerId)
	query += " RETURNING *;"

	var updatedPassenger models.Passenger
	if err := s.db.QueryRowxContext(ctx, query, params...).StructScan(&updatedPassenger); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &updatedPassenger, nil
}

func (s *StorageRepo) UpdateDocumentInfo(ctx context.Context, documentId string, updateData requests.DocumentUpdateRequest) (*models.Document, error) {
	const op = "storage.postgresql.UpdateDocumentInfo"

	query := "UPDATE Document SET"
	params := []interface{}{}
	idx := 1
	if updateData.PassengerId != nil {
		if idx > 1 {
			query += ", "
		}
		query += " passenger_id = $" + strconv.Itoa(idx)
		params = append(params, *updateData.PassengerId)
		idx++
	}
	if updateData.DocumentType != nil {
		if idx > 1 {
			query += ", "
		}
		query += " document_type = $" + strconv.Itoa(idx)
		params = append(params, *updateData.DocumentType)
		idx++
	}
	if updateData.DocumentNumber != nil {
		if idx > 1 {
			query += ", "
		}
		query += " document_number = $" + strconv.Itoa(idx)
		params = append(params, *updateData.DocumentNumber)
		idx++
	}

	query += " WHERE document_id = $" + strconv.Itoa(idx)
	idx++
	params = append(params, documentId)
	query += " RETURNING document_id, passenger_id AS \"document.passenger_id\", document_type, document_number;"

	var updatedDocument models.Document
	if err := s.db.QueryRowxContext(ctx, query, params...).StructScan(&updatedDocument); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &updatedDocument, nil
}

func (s *StorageRepo) DeleteTicketById(ctx context.Context, ticketId string) error {
	const op = "storage.postgres.DeleteTicketById"

	_, err := s.db.QueryxContext(ctx, query.DeleteTicketPassenger, ticketId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.db.QueryxContext(ctx, query.DeleteTicket, ticketId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *StorageRepo) DeletePassengerById(ctx context.Context, passengerId string) error {
	const op = "storage.postgres.DeletePassengerById"

	_, err := s.db.QueryxContext(ctx, query.DeletePassenger, passengerId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *StorageRepo) DeleteDocumentById(ctx context.Context, documentId string) error {
	const op = "storage.postgres.DeleteDocumentById"

	_, err := s.db.QueryxContext(ctx, query.DeleteDocument, documentId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *StorageRepo) GetPassengerReport(ctx context.Context, passengerId string, startDate time.Time, endDate time.Time) ([]*response.FlightReport, error) {
	const op = "storage.postgres.GetPassengerReport"

	var report []*response.FlightReport
	err := s.db.SelectContext(ctx, &report, query.GetPassengerReport, startDate, endDate, passengerId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return report, nil
}
