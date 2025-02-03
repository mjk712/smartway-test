package test_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"log/slog"
	"net/http"
	"smartway-test/internal/config"
	http_server "smartway-test/internal/http-server"
	"smartway-test/internal/http-server/requests"
	"smartway-test/internal/http-server/response"
	"smartway-test/internal/models"
	"smartway-test/internal/service"
	"smartway-test/internal/storage"
	"smartway-test/internal/test/testdata"
	"smartway-test/internal/test/testhelpers"
	"testing"
	"time"
)

const (
	ticketNumber = "124237694"
	ticketId     = "1"
	passengerID  = "1"
	documentId   = "1"
	startDate    = "2020-11-20"
	endDate      = "2030-12-26"
)

type HandlersTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	service     service.FlightService
	db          *sqlx.DB
	repository  *storage.StorageRepo
	ctx         context.Context
	server      *http.Server
	config      *config.Config
}

func (suite *HandlersTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	suite.config = &config.Config{HTTPServer: config.HTTPServer{Address: "127.0.0.1:8081"}}

	//запуск контейнера с тестовой postgres
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	suite.Require().NoError(err)
	suite.pgContainer = pgContainer

	db, err := sqlx.Connect("postgres", suite.pgContainer.ConnectionString)
	suite.Require().NoError(err)

	suite.db = db

	//запуск миграций
	suite.RunMigrations()
	//очистка и заполнение тестовых данных
	suite.ClearData()
	suite.FillData()

	//запуск сервиса и сервера
	suite.repository = &storage.StorageRepo{DB: suite.db}
	suite.service = service.NewFlightService(suite.repository)
	suite.server = http_server.NewServer(suite.ctx, slog.Default(), suite.config, suite.service)

	go func() {
		err := suite.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			suite.T().Fatal(err)
		}
	}()
	time.Sleep(1 * time.Second)
}

func (suite *HandlersTestSuite) TearDownSuite() {
	if err := suite.server.Shutdown(suite.ctx); err != nil {
		suite.T().Fatal(err)
	}

	err := suite.pgContainer.Terminate(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *HandlersTestSuite) RunMigrations() {
	// выполняем миграции
	m, err := migrate.New("file://../storage/migrations", suite.pgContainer.ConnectionString)
	if err != nil {
		log.Fatalf("failed to create migration instance: %s", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %s", err)
	}
}

func (suite *HandlersTestSuite) FillData() {
	_, err := suite.db.QueryxContext(suite.ctx, testdata.InsertTestData)
	if err != nil {
		log.Fatalf("failed to insert test data: %s", err)
	}
}

func (suite *HandlersTestSuite) ClearData() {
	_, err := suite.db.QueryxContext(suite.ctx, testdata.ClearTestData)
	if err != nil {
		log.Fatalf("failed to clear test data: %s", err)
	}
}

func (suite *HandlersTestSuite) TestGetTicketsHandler() {

	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8081/api/tickets", nil)
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var tickets []models.Ticket
	err = json.NewDecoder(resp.Body).Decode(&tickets)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tickets)

}

func (suite *HandlersTestSuite) TestGetTicketFullInfoHandler() {

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8081/api/ticket/"+ticketNumber, nil)
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var fullTicketInfo response.FullTicketInfo
	err = json.NewDecoder(resp.Body).Decode(&fullTicketInfo)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), ticketNumber, fullTicketInfo.Ticket.OrderNumber)
	assert.NotEmpty(suite.T(), fullTicketInfo.Passengers)
}

func (suite *HandlersTestSuite) TestGetPassengersByTicketNumberHandler() {
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8081/api/passengers/"+ticketNumber, nil)
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var passengers []models.Passenger
	err = json.NewDecoder(resp.Body).Decode(&passengers)
	assert.NoError(suite.T(), err)

	assert.NotEmpty(suite.T(), passengers)
	for _, passenger := range passengers {
		assert.NotEmpty(suite.T(), passenger)
	}
}

func (suite *HandlersTestSuite) TestGetPassengerReportHandler() {
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8081/api/reports/passenger/"+passengerID, nil)
	q := req.URL.Query()
	q.Add("start_date", startDate)
	q.Add("end_date", endDate)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var report []response.FlightReport
	err = json.NewDecoder(resp.Body).Decode(&report)
	assert.NoError(suite.T(), err)

	assert.NotEmpty(suite.T(), report)
	for _, flight := range report {
		assert.NotEmpty(suite.T(), flight)
	}
}

func (suite *HandlersTestSuite) TestGetDocumentsByPassengerIDHandler() {
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8081/api/documents/"+passengerID, nil)
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var documents []models.Document
	err = json.NewDecoder(resp.Body).Decode(&documents)
	assert.NoError(suite.T(), err)

	assert.NotEmpty(suite.T(), documents)
	for _, document := range documents {
		assert.NotEmpty(suite.T(), document)
	}
}

func (suite *HandlersTestSuite) TestUpdateTicketInfoHandler() {
	newDestination := "Texas"
	updateRequest := requests.TicketUpdateRequest{
		DestinationPoint: &newDestination,
	}

	requestBody, err := json.Marshal(updateRequest)
	assert.NoError(suite.T(), err)

	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8081/api/ticket/2", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var ticket models.Ticket
	err = json.NewDecoder(resp.Body).Decode(&ticket)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), newDestination, ticket.DestinationPoint)
}

func (suite *HandlersTestSuite) TestUpdatePassengerInfoHandler() {
	newFirstName := "John"
	newLastName := "Doe"
	updateRequest := requests.UpdatePassengerRequest{
		FirstName: &newFirstName,
		LastName:  &newLastName,
	}

	requestBody, err := json.Marshal(updateRequest)
	assert.NoError(suite.T(), err)

	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8081/api/passenger/"+passengerID, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var passenger models.Passenger
	err = json.NewDecoder(resp.Body).Decode(&passenger)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), newFirstName, passenger.FirstName)
	assert.Equal(suite.T(), newLastName, passenger.LastName)

}

func (suite *HandlersTestSuite) TestUpdateDocumentInfoHandler() {
	newDocumentType := "Studak"
	newDocumentNumber := "12345"

	updatedRequest := requests.DocumentUpdateRequest{
		DocumentType:   &newDocumentType,
		DocumentNumber: &newDocumentNumber,
	}

	requestBody, err := json.Marshal(updatedRequest)
	assert.NoError(suite.T(), err)

	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8081/api/document/"+documentId, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var document models.Document
	err = json.NewDecoder(resp.Body).Decode(&document)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), newDocumentType, document.DocumentType)
	assert.Equal(suite.T(), newDocumentNumber, document.DocumentNumber)
}

func (suite *HandlersTestSuite) TestDeleteTicketHandler() {
	req, _ := http.NewRequest(http.MethodDelete, "http://127.0.0.1:8081/api/ticket/"+ticketId, nil)
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var responseMessage string
	err = json.NewDecoder(resp.Body).Decode(&responseMessage)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), fmt.Sprintf("ticket with id %s deleted", ticketId), responseMessage)

	var ticket models.Ticket
	err = suite.db.Get(&ticket, "SELECT * FROM flight_ticket WHERE id = $1", ticketId)
	assert.Error(suite.T(), err)
}

func (suite *HandlersTestSuite) TestDeletePassengerHandler() {
	passId := "2"

	req, _ := http.NewRequest(http.MethodDelete, "http://127.0.0.1:8081/api/passenger/"+passId, nil)
	client := &http.Client{}
	resp, err := client.Do(req)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var responseMessage string
	err = json.NewDecoder(resp.Body).Decode(&responseMessage)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), fmt.Sprintf("passenger with id %s deleted", passId), responseMessage)

	var passenger models.Passenger
	err = suite.db.Get(&passenger, "SELECT * FROM passenger WHERE id = $1", passId)
	assert.Error(suite.T(), err)
}

// Запуск тестов
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
