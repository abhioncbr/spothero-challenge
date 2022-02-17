package model

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

// Suite has gorm DB, sql DB  and db instances
type Suite struct {
	suite.Suite
	DB   *gorm.DB
	sqlDB *sql.DB
	mock sqlmock.Sqlmock
	rate Rate
}

func (s *Suite) SetupSuite() {
	mock, DB, sqlDB := GetDatabase(s)
	s.mock = mock
	s.DB = DB
	s.sqlDB = sqlDB

	s.rate = Rate{
		Days: "mon,tues,thurs",
		Times: "0900-2100",
		Tz:  "America/Chicago",
		Price: 1500,
	}
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestDBMigrate(){
	s.mock.ExpectExec("CREATE TABLE `rates`(.*)").WillReturnResult(sqlmock.NewResult(0, 1))
	dbMigrateError := DBMigrate(s.DB)
	require.NoError(s.T(), dbMigrateError)
}

func (s *Suite) TestLoadRatesOnStart() {
	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT INTO `rates`(.*)").
		WithArgs(s.rate.Days, s.rate.Times, s.rate.Tz, s.rate.Price).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	loadError := LoadRatesOnStart("mock_rate.json", s.DB)
	require.NoError(s.T(), loadError)
}

// TestLoadRatesOnStartWithLoadedData: test the LoadRatesOnStart with already loaded data.
func (s *Suite) TestLoadRatesOnStartWithLoadedData() {
	rows := s.mock.NewRows([]string{"days", "times", "tz", "price"}).AddRow(s.rate.Days, s.rate.Times, s.rate.Tz, s.rate.Price)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rates`")).WillReturnRows(rows)

	loadError := LoadRatesOnStart("mock_rate.json", s.DB)
	require.NoError(s.T(), loadError)
}

// GetDatabase: set the sql mock and gorm v=based DB for testing.
func GetDatabase(s *Suite) (sqlmock.Sqlmock, *gorm.DB, *sql.DB){
	sqlDB, mock, err := sqlmock.NewWithDSN("sql_mock_db", sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(s.T(), err)
	require.NotNil(s.T(), sqlDB)
	require.NotNil(s.T(), mock)

	// A select sqlite_version() query will be run when gorm opens the database
	// we need to expect that here
	columns := []string{"version"}
	mock.ExpectQuery("select sqlite_version()").WithArgs().WillReturnRows(
		mock.NewRows(columns).FromCSVString("1"),
	)
	DB, err := gorm.Open(sqlite.Dialector{DSN: "sql_mock_db", Conn: sqlDB}, &gorm.Config{})
	require.NotNil(s.T(), DB)
	require.NoError(s.T(), err)

	return mock, DB, sqlDB
}