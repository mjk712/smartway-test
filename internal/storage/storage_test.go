package storage

import (
	"context"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"smartway-test/internal/http-server/requests"
	"smartway-test/internal/storage/testdata"
	"testing"
	"time"
)

var db *sqlx.DB

func setupTestDB(t *testing.T) *postgres.PostgresContainer {

	ctx := context.Background()

	dbName := "testdb"
	dbUser := "testuser"
	dbPassword := "testpass"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:14",
		//postgres.WithInitScripts(filepath.Join("testdata", "init-user-db.sh")),
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
		t.Fatalf("failed to start test db container: %s", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %s", err)
	}
	//подключаемся в бд
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatalf("failed to connect to test db: %s", err)
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

	storage := &StorageRepo{db: db}

	//готовим тестовые данные
	_, err := db.QueryxContext(ctx, testdata.InsertTestData)
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

	storage := &StorageRepo{db: db}

	//готовим тестовые данные
	_, err := db.QueryxContext(ctx, testdata.InsertTestData)
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
