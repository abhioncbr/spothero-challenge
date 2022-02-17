package handler

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"regexp"
)

// TestStandardHandleRateTimes test the start and end time as per correct format.
func (s *Suite) TestStandardHandleRateTimes(){
	start, end, err := handleRateTimes(s.rate)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *start, 9)
	assert.Equal(s.T(), *end, 21)
}

// TestNotStandardHandleRateTimes should throw error for incorrect format.
func (s *Suite) TestNotStandardHandleRateTimes(){
	s.rate.Times = "09002100"
	start, end, err := handleRateTimes(s.rate)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "time value is not as per the standard")
	assert.Nil(s.T(), start)
	assert.Nil(s.T(), end)
	s.rate.Times = "0900-2100"
}

// TestNonParsableStartHandleRateTimes should throw error for non-parsable start hour.
func (s *Suite) TestNonParsableStartHandleRateTimes(){
	s.rate.Times = "0a00-2100"
	start, end, err := handleRateTimes(s.rate)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "strconv.Atoi: parsing \"a\": invalid syntax")
	assert.Nil(s.T(), start)
	assert.Nil(s.T(), end)
	s.rate.Times = "0900-2100"
}

// TestNonParsableEndHandleRateTimes should throw error for non-parsable start hour.
func (s *Suite) TestNonParsableEndHandleRateTimes(){
	s.rate.Times = "0900-2Z00"
	start, end, err := handleRateTimes(s.rate)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), err.Error(), "strconv.Atoi: parsing \"2Z\": invalid syntax")
	assert.Nil(s.T(), start)
	assert.Nil(s.T(), end)
	s.rate.Times = "0900-2100"
}

// TestValidateTimeStartParam should return the valid 'start' time object
func (s *Suite) TestValidateTimeStartParam(){
	req, err := http.NewRequest("GET", "/price?start=2015-07-01T07:00:00-05:00", nil)
	assert.NoError(s.T(), err)
	startTime, validationErr := validateTimeParam(req.URL, "start")
	assert.NoError(s.T(), validationErr)
	assert.NotNil(s.T(), startTime)
	assert.Equal(s.T(), startTime.Hour(), 7)
}

// TestValidateTimeEndParam should return the valid 'end' time object
func (s *Suite) TestValidateTimeEndParam(){
	req, err := http.NewRequest("GET", "/price?end=2015-07-01T07:00:00-05:00", nil)
	assert.NoError(s.T(), err)
	endTime, validationErr := validateTimeParam(req.URL, "end")
	assert.NoError(s.T(), validationErr)
	assert.NotNil(s.T(), endTime)
	assert.Equal(s.T(), endTime.Hour(), 7)
}

// TestValidateTimeMissingParam test should throw error for missing param
func (s *Suite) TestValidateTimeMissingParam(){
	req, err := http.NewRequest("GET", "/price", nil)
	assert.NoError(s.T(), err)
	time, validationErr := validateTimeParam(req.URL, "start")
	assert.Error(s.T(), validationErr)
	assert.Equal(s.T(), validationErr.Error(), "missing Url Param 'start' ")
	assert.Nil(s.T(), time)
}

// TestValidateTimeMissingParamValue test should throw error for missing param value
func (s *Suite) TestValidateTimeMissingParamValue(){
	req, err := http.NewRequest("GET", "/price?start", nil)
	assert.NoError(s.T(), err)
	time, validationErr := validateTimeParam(req.URL, "start")
	assert.Error(s.T(), validationErr)
	assert.Equal(s.T(), validationErr.Error(), "Url param 'start' has no value ")
	assert.Nil(s.T(), time)
}

// TestValidateTimeNotParsableIso8601Time test should throw error for not parsable time value
func (s *Suite) TestValidateTimeNotParsableIso8601Time(){
	req, err := http.NewRequest("GET", "/price?start=qq", nil)
	assert.NoError(s.T(), err)
	time, validationErr := validateTimeParam(req.URL, "start")
	assert.Error(s.T(), validationErr)
	assert.Equal(s.T(), validationErr.Error(), "Url param 'start' isn't as per ISO-8601 standard ")
	assert.Nil(s.T(), time)
}

// TestGetPriceValidPrice1500 return valid 1500 price for the query.
func (s *Suite) TestGetPriceValidPrice1500(){
	rows := s.mock.NewRows([]string{"days", "times", "tz", "price"}).AddRow(s.rate.Days, s.rate.Times, s.rate.Tz, s.rate.Price)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rates`")).WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/price?start=2015-07-04T15:00:00-05:00&end=2015-07-04T20:00:00-05:00", nil)
	assert.NoError(s.T(), err)
	httpRec := httptest.NewRecorder()

	GetPrice(s.DB, httpRec, req)
	assert.Equal(s.T(), httpRec.Code, http.StatusOK)

	price := Price{Price: 1500}
	jsonPrice, marshalError := json.Marshal(price)
	assert.NoError(s.T(), marshalError)
	assert.Equal(s.T(), httpRec.Body.String(), string(jsonPrice) )
}

// TestGetPriceValidPrice1750 return valid 1750 price for the query.
func (s *Suite) TestGetPriceValidPrice1750(){
	rows := s.mock.NewRows([]string{"days", "times", "tz", "price"}).
		AddRow(s.rate.Days, s.rate.Times, s.rate.Tz, s.rate.Price).
		AddRow("wed", "0600-1800", "America/Chicago", 1750)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rates`")).WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/price?start=2015-07-01T07:00:00-05:00&end=2015-07-01T12:00:00-05:00", nil)
	assert.NoError(s.T(), err)
	httpRec := httptest.NewRecorder()

	GetPrice(s.DB, httpRec, req)
	assert.Equal(s.T(), httpRec.Code, http.StatusOK)

	price := Price{Price: 1750}
	jsonPrice, marshalError := json.Marshal(price)
	assert.NoError(s.T(), marshalError)
	assert.Equal(s.T(), httpRec.Body.String(), string(jsonPrice) )
}

// TestGetPriceUnavailable return Unavailable price for the query.
func (s *Suite) TestGetPriceUnavailable(){
	rows := s.mock.NewRows([]string{"days", "times", "tz", "price"}).
		AddRow(s.rate.Days, s.rate.Times, s.rate.Tz, s.rate.Price).
		AddRow("wed", "0600-1800", "America/Chicago", 1750)
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `rates`")).WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/price?start=2015-07-04T07:00:00%2B05:00&end=2015-07-04T20:00:00%2B05:00", nil)
	assert.NoError(s.T(), err)
	httpRec := httptest.NewRecorder()

	GetPrice(s.DB, httpRec, req)
	assert.Equal(s.T(), httpRec.Code, http.StatusOK)
	assert.Equal(s.T(), httpRec.Body.String(), "\"unavailable\"" )
}
