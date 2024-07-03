package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type mockDB struct {
	dbConn  *gorm.DB
	sqlmock sqlmock.Sqlmock
}

func newMockDB(t *testing.T) mockDB {
	sqlDB, sqlmock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})
	dbConn, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return mockDB{
		dbConn:  dbConn,
		sqlmock: sqlmock,
	}
}

func TestRunMigrations(t *testing.T) {
	mDB := newMockDB(t)
	db, err := mDB.dbConn.DB()
	assert.NoError(t, err)
	defer db.Close()

	mDB.sqlmock.ExpectQuery(`^SELECT count\(\*\) FROM information_schema\.tables WHERE table_schema = CURRENT_SCHEMA\(\) AND table_name = \$1 AND table_type = \$2$`).
		WithArgs("blocks", "BASE TABLE").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mDB.sqlmock.ExpectExec(`^CREATE TABLE "blocks"`).WillReturnResult(sqlmock.NewResult(0, 1))
	mDB.sqlmock.ExpectExec(`^CREATE INDEX IF NOT EXISTS "idx_blocks_number" ON "blocks" \("number" asc\)`).WillReturnResult(sqlmock.NewResult(0, 1))

	err = runMigrations(mDB.dbConn)
	assert.NoError(t, err)

	err = mDB.sqlmock.ExpectationsWereMet()
	assert.NoError(t, err)
}
