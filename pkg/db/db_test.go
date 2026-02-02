package db

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type suite struct {
	dbMock  DB
	sqlMock sqlmock.Sqlmock
}

func newSuite(t *testing.T) suite {
	sqlDB, sqlMock, err := sqlmock.New()
	assert.NoError(t, err)
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)
	return suite{
		dbMock:  &GormDB{gormDB},
		sqlMock: sqlMock,
	}
}

func TestRunMigrations(t *testing.T) {
	sqlDB, sqlMock, err := sqlmock.New()
	assert.NoError(t, err)
	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	sqlMock.ExpectQuery(`^SELECT count\(\*\) FROM information_schema\.tables WHERE table_schema = CURRENT_SCHEMA\(\) AND table_name = \$1 AND table_type = \$2$`).
		WithArgs("blocks", "BASE TABLE").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	sqlMock.ExpectExec(`^CREATE TABLE "blocks"`).WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectExec(`^CREATE INDEX IF NOT EXISTS "idx_blocks_number" ON "blocks" \("number" asc\)`).WillReturnResult(sqlmock.NewResult(0, 1))
	sqlMock.ExpectExec(`^CREATE TABLE "transactions" \("hash" char\(66\),"from" char\(42\) NOT NULL,"to" char\(42\),"contract" char\(66\) NOT NULL,"value" numeric NOT NULL,"data" bytea NOT NULL,"gas" numeric NOT NULL,"gas_price" numeric NOT NULL,"cost" numeric NOT NULL,"nonce" numeric NOT NULL,"status" numeric NOT NULL,"block_hash" char\(66\) NOT NULL,PRIMARY KEY \("hash"\)\)$`).WillReturnResult(sqlmock.NewResult(0, 1))

	err = runMigrations(gormDB)
	assert.NoError(t, err)

	err = sqlMock.ExpectationsWereMet()
	assert.NoError(t, err)
}
