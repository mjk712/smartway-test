package test

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"smartway-test/internal/http-server/requests"
	"smartway-test/internal/storage"
	"smartway-test/internal/test/testdata"
	"smartway-test/internal/test/testhelpers"
	"testing"
)

type StorageRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repository  *storage.StorageRepo
	ctx         context.Context
	db          *sqlx.DB
}

func (suite *StorageRepoTestSuite) FillData() {
	_, err := suite.repository.DB.QueryxContext(suite.ctx, testdata.InsertTestData)
	if err != nil {
		log.Fatalf("failed to insert test data: %s", err)
	}
}

func (suite *StorageRepoTestSuite) ClearData() {
	_, err := suite.repository.DB.QueryxContext(suite.ctx, testdata.ClearTestData)
	if err != nil {
		log.Fatalf("failed to clear test data: %s", err)
	}
}

func (suite *StorageRepoTestSuite) RunMigrations() {
	// выполняем миграции
	m, err := migrate.New("file://../storage/migrations", suite.pgContainer.ConnectionString)
	if err != nil {
		log.Fatalf("failed to create migration instance: %s", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrations: %s", err)
	}
}

func (suite *StorageRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	suite.Require().NoError(err)

	suite.pgContainer = pgContainer

	suite.db, err = sqlx.Connect("postgres", suite.pgContainer.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	suite.RunMigrations()

	suite.repository = &storage.StorageRepo{DB: suite.db}

	suite.ClearData()

	suite.FillData()
}

func (suite *StorageRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("failed to terminate postgres container: %s", err)
	}
}

func (suite *StorageRepoTestSuite) TestStorageRepo_GetTickets() {
	t := suite.T()

	tickets, err := suite.repository.GetTickets(suite.ctx)
	assert.NoError(t, err)
	assert.NotNil(t, tickets)
	assert.Len(t, tickets, 6)
	assert.Equal(t, "Penza", tickets[0].DestinationPoint)
}

func (suite *StorageRepoTestSuite) TestStorageRepo_GetTicket() {
	t := suite.T()
	testDestinationPoint := "Sever"

	//вызываем тестируемый метод
	updateRequest := requests.TicketUpdateRequest{
		DestinationPoint: &testDestinationPoint,
	}
	updatedTicket, err := suite.repository.UpdateTicketInfo(suite.ctx, "2", updateRequest)
	assert.NoError(t, err)
	assert.Equal(t, "Sever", updatedTicket.DestinationPoint)
}

func (suite *StorageRepoTestSuite) TestStorageRepo_GetPassengersByTicketNumber() {
	t := suite.T()

	ticketNumber := "124237694"

	passengers, err := suite.repository.GetPassengersByTicketNumber(suite.ctx, ticketNumber)
	assert.NoError(t, err)
	assert.NotNil(t, passengers)
	assert.Equal(t, "Valerich", passengers[0].FirstName)
	assert.Equal(t, "Matukov", passengers[0].LastName)
}

func (suite *StorageRepoTestSuite) TestStorageRepo_GetDocumentsByPassengerId() {
	t := suite.T()

	passengerID := "1"

	// Вызываем тестируемый метод
	documents, err := suite.repository.GetDocumentsByPassengerId(suite.ctx, passengerID)
	assert.NoError(t, err)
	assert.NotNil(t, documents)
	assert.Len(t, documents, 3)
	assert.Equal(t, "passport", documents[0].DocumentType)
	assert.Equal(t, "34 25 876527", documents[0].DocumentNumber)
}

func TestStorageRepoTestSuite(t *testing.T) {
	suite.Run(t, new(StorageRepoTestSuite))
}
