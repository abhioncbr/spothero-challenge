package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
)

// Rate struct for storing the rate properties in DB
type Rate struct {
	Days string `gorm:"primaryKey" json:"days"`
	Times string `gorm:"primaryKey" json:"times"`
	Tz    string `gorm:"primaryKey" json:"tz"`
	Price int    `json:"price"`
}

// Rates struct contains the list of rate.
type Rates struct {
	Rates []Rate `json:"rates"`
}

// DBMigrate migrate the DB on app start and registering the models(rate)
func DBMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Rate{})
}

// LoadRatesOnStart save the provided rate data in DB if not already present
func LoadRatesOnStart(ratesFile string, db *gorm.DB) error {
	var obRates []Rate
	db.Find(&obRates)
	if len(obRates) == 0{
		file, fileErr := ioutil.ReadFile(ratesFile)

		// if error in reading the rate list json file, return error
		if fileErr != nil{
			return fileErr
		}

		var rates Rates
		ratesLoadErr := json.Unmarshal([]byte(file), &rates)

		// if error in unmarshalling rate list json, return error
		if ratesLoadErr != nil{
			return ratesLoadErr
		}

		// saving initial rate list in the DB
		return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rates.Rates).Error
	}
	return nil
}