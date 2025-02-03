package test_test

import (
	"context"
	"encoding/json"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"log/slog"
	"net/http"
	"smartway-test/internal/config"
	http_server "smartway-test/internal/http-server"
	"smartway-test/internal/models"
	"smartway-test/internal/service"
	"smartway-test/internal/storage"
	"smartway-test/internal/test/testdata"
	"smartway-test/internal/test/testhelpers"
	"testing"
	"time"
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

// Запуск тестов
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
