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

func TestStorageRepoTestSuite(t *testing.T) {
	suite.Run(t, new(StorageRepoTestSuite))
}

/*func setupTestDB(t *testing.T) *postgres.PostgresContainer {

	ctx := context.Background()

	dbName := "testdb"
	dbUser := "testuser"
	dbPassword := "testpass"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:14",
		//postgres.WithInitScripts(filepath.Join("testdata", "init-user-DB.sh")),
		//postgres.WithConfigFile(filepath.Join("testdata", "my-postgres.conf")),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		t.Fatalf("failed to start test DB container: %s", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}
	//подключаемся в бд
	DB, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect to test DB: %s", err)
	}

	// выполняем миграции
	m, err := migrate.New("file://migrations", connStr)
	if err != nil {
		t.Fatalf("failed to create migration instance: %s", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatalf("failed to run migrations: %s", err)
	}

	return postgresContainer

}

func TestMain(m *testing.M) {
	//TODO Используем для тестов один контейнер, миграции делаем один раз
}

func TestStorageRepo_GetTickets(t *testing.T) {
	ctx := context.Background()
	postgresContainer := setupTestDB(t)

	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	storage := &StorageRepo{DB: DB}

	//готовим тестовые данные
	_, err := DB.QueryxContext(ctx, testdata.InsertTestData)
	if err != nil {
		t.Fatalf("failed to insert test data: %s", err)
	}

	//вызываем тестируемый метод
	tickets, err := storage.GetTickets(ctx)
	assert.NoError(t, err)
	assert.Len(t, tickets, 6)
	assert.Equal(t, "Penza", tickets[0].DestinationPoint)
}

func TestStorageRepo_UpdateTicketInfo(t *testing.T) {
	ctx := context.Background()
	postgresContainer := setupTestDB(t)

	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	storage := &StorageRepo{DB: DB}

	//готовим тестовые данные
	_, err := DB.QueryxContext(ctx, testdata.InsertTestData)
	if err != nil {
		t.Fatalf("failed to insert test data: %s", err)
	}

	testDestinationPoint := "Sever"

	//вызываем тестируемый метод
	updateRequest := requests.TicketUpdateRequest{
		DestinationPoint: &testDestinationPoint,
	}
	updatedTicket, err := storage.UpdateTicketInfo(ctx, "2", updateRequest)
	assert.NoError(t, err)
	assert.Equal(t, "Sever", updatedTicket.DestinationPoint)
}
*/
