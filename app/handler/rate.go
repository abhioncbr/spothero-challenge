package handler

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"net/http"
	"spotHero/app/model"
)

// GetAllRates api endpoints to get all the rates stored in the database.
func GetAllRates(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var rates []model.Rate
	getError := db.Find(&rates).Error
	if getError != nil {
		respondError(w, http.StatusBadRequest, getError.Error())
	}
	respondJSON(w, http.StatusOK, rates)
}

// PutRate api endpoints to upsert the rate in the database
func PutRate(db *gorm.DB, w http.ResponseWriter, r *http.Request){
	rate := model.Rate{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rate); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	upsertErr := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "days"}, {Name: "times"}, {Name: "tz"}}, // key colume
		DoUpdates: clause.AssignmentColumns([]string{"price"}), // column needed to be updated
	}).Create(&rate).Error

	if upsertErr != nil {
		respondError(w, http.StatusInternalServerError, upsertErr.Error())
		return
	}

	respondJSON(w, http.StatusCreated, rate)
}
