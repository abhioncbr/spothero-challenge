package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"regexp"
	"spotHero/app/model"
	"testing"
)

// Suite has gorm DB, sql DB  and db instances
type Suite struct {
	suite.Suite
	DB   *gorm.DB
	sqlDB *sql.DB
	mock sqlmock.Sqlmock
	rate model.Rate
}

// SetupSuite setting the suite for testing.
func (s *Suite) SetupSuite() {
	mock, DB, sqlDB := GetDatabase(s)
	s.mock = mock
	s.DB = DB
	s.sqlDB = sqlDB

	s.rate = model.Rate{
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

// TestGetAllRatesHandler test the GetAllRates endpoint.
func (s *Suite) TestGetAllRatesHandler() {
	rows := s.mock.NewRows([]string{"days", "times", "tz", "price"}).AddRow(s.rate.Days, s.rate.Times, s.rate.Tz, s.rate.Price)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rates`")).WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/rates", nil)
	assert.NoError(s.T(), err)
	httpRec := httptest.NewRecorder()

	GetAllRates(s.DB, httpRec, req)
	assert.Equal(s.T(), httpRec.Code, http.StatusOK)

	var rates []model.Rate
	rates = append(rates, s.rate)
	jsonRates, marshalError := json.Marshal(rates)
	assert.NoError(s.T(), marshalError)
	assert.Equal(s.T(), httpRec.Body.String(), string(jsonRates) )
}

// TestPutRateInsert test to update the already stored rate
func (s *Suite) TestPutRateInsert(){
	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT INTO `rates`(.*)").
		WithArgs(s.rate.Days, s.rate.Times, s.rate.Tz, s.rate.Price).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	jsonRate, marshalError := json.Marshal(s.rate)
	assert.NoError(s.T(), marshalError)


	req, err := http.NewRequest("PUT", "/rates", bytes.NewBuffer(jsonRate))
	assert.NoError(s.T(), err)
	httpRec := httptest.NewRecorder()
	PutRate(s.DB, httpRec, req)
	assert.Equal(s.T(), httpRec.Code, http.StatusCreated)
	assert.Equal(s.T(), httpRec.Body.String(), string(jsonRate) )
}

func (s *Suite) TestPutRateUpdate(){
	s.rate.Price = 4000
	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT INTO `rates`(.*)").
		WithArgs(s.rate.Days, s.rate.Times, s.rate.Tz, s.rate.Price).
		WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()

	jsonRate, marshalError := json.Marshal(s.rate)
	assert.NoError(s.T(), marshalError)

	req, err := http.NewRequest("PUT", "/rates", bytes.NewBuffer(jsonRate))
	assert.NoError(s.T(), err)
	httpRec := httptest.NewRecorder()
	PutRate(s.DB, httpRec, req)
	assert.Equal(s.T(), httpRec.Code, http.StatusCreated)
	assert.Equal(s.T(), httpRec.Body.String(), string(jsonRate) )
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