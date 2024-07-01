package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type mockDB struct {
	dbConn  *gorm.DB
	sqlmock sqlmock.Sqlmock
}

func newMockDB(t *testing.T) mockDB {
	sqlDB, sqlmock, err := sqlmock.New()
	if err != nil {
		t.Errorf("failed to create sqlmock %v", err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})
	dbConn, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Errorf("failed to open gorm v2 db, got error: %v", err)
	}

	return mockDB{
		dbConn:  dbConn,
		sqlmock: sqlmock,
	}
}

func TestRunMigrations(t *testing.T) {
	mDB := newMockDB(t)
	db, err := mDB.dbConn.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mDB.sqlmock.ExpectQuery(`^SELECT count\(\*\) FROM information_schema\.tables WHERE table_schema = CURRENT_SCHEMA\(\) AND table_name = \$1 AND table_type = \$2$`).
		WithArgs("blocks", "BASE TABLE").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mDB.sqlmock.ExpectExec(`^CREATE TABLE "blocks"`).WillReturnResult(sqlmock.NewResult(0, 1))
	mDB.sqlmock.ExpectExec(`^CREATE INDEX IF NOT EXISTS "idx_blocks_number" ON "blocks" \("number" asc\)`).WillReturnResult(sqlmock.NewResult(0, 1))

	err = runMigrations(mDB.dbConn)
	if err != nil {
		t.Fatalf("Error running migrations: %v", err)
	}
	if err := mDB.sqlmock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
