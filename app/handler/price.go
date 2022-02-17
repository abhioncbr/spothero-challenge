package handler

import (
	"errors"
	"fmt"
	"github.com/relvacode/iso8601"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"spotHero/app/model"
	"strconv"
	"strings"
	"time"
)

// Price contains the price for response.
type Price struct {
	Price int `json:"price"`
}

// GetPrice return the price based on query start and end time param
func GetPrice(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	startTime, startErr := validateTimeParam(r.URL, "start")
	endTime, endErr := validateTimeParam(r.URL, "end")

	if startErr != nil {
		respondError(w, http.StatusOK, startErr.Error())
		return
	}

	if endErr != nil {
		respondError(w, http.StatusOK, endErr.Error())
		return
	}

	timeDifference := int(endTime.Sub(*startTime).Hours())
	if timeDifference > 24 || timeDifference < 0 {
		respondJSON(w, http.StatusOK, "unavailable")
		return
	}

	loc, _ := time.LoadLocation("America/Chicago")

	// getting the rates from the database
	var obRates []model.Rate
	dayRune:= []rune(startTime.Weekday().String())
	day := "%" + strings.ToLower(string(dayRune[0:2])) +"%"
	db.Where("days like ? AND tz =?", day, loc.String()).Find(&obRates)

	// finding the correct price as per the stored rates.
	for _, rate := range obRates {
		rStartTime, rEndTime, err := handleRateTimes(rate)
		if err != nil {
			respondJSON(w, http.StatusOK, "unavailable")
		}

		if startTime.Hour() >= *rStartTime && endTime.Hour() <= *rEndTime {
			respondJSON(w, http.StatusOK, Price{rate.Price})
			return
		}
	}

	respondJSON(w, http.StatusOK, "unavailable")
}

// validateTimeParam validate the time param from the http request query
func validateTimeParam(url *url.URL, paramName string) (*time.Time, error){
	param, isPresent := url.Query()[paramName]
	// check if paramName is present in the URL query or not, if not return error
	if !isPresent {
		return nil, errors.New(fmt.Sprintf("missing Url Param '%s' ", paramName))
	}

	// check if the length of the value of param is not zero
	if len(param) < 0 || param[0] == "" {
		return nil, errors.New(fmt.Sprintf("Url param '%s' has no value ", paramName))
	}

	// parsed the time param based on ISO-8601 standard.
	parsedTime, parsErr := iso8601.ParseString(param[0])
	if parsErr != nil{
		return nil, errors.New(fmt.Sprintf("Url param '%s' isn't as per ISO-8601 standard ", paramName))
	}

	return &parsedTime, parsErr
}

// handleRateTimes handle the time string value of the rate model for processing.
func handleRateTimes(rate model.Rate) (*int, *int, error) {
	times := rate.Times
	splitTime := strings.SplitAfter(times,"-") // format: "start-end", split time strings into two parts

	// after split, it should be into two parts.
	if len(splitTime) != 2 {
		return nil, nil, errors.New("time value is not as per the standard")
	}

	// trim the '0' and parse the value into int for start hour
	startTimeString := strings.Trim(strings.TrimSuffix(splitTime[0], "-"), "0")
	parsedStartTime, err := strconv.Atoi(startTimeString)
	if err != nil {
		return nil, nil, err
	}

	// trim the '0' and parse the value into int for end hour
	endTimeString := strings.Trim(splitTime[1], "0")
	parsedEndTime, err := strconv.Atoi(endTimeString)
	if err != nil{
		return nil, nil, err
	}

	return &parsedStartTime, &parsedEndTime, nil
}