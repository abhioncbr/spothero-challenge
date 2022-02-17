package main

import (
	"fmt"
	"spotHero/app"
	"spotHero/config"
	_ "time/tzdata"
)

const(
	// AppPort default port of the app
	AppPort = 5000
)

// main method of the app
func main() {
	spotHeroApp := &app.App{}
	dbConfig := config.GetSqliteConfig("./rates.db")
	spotHeroApp.Initialize(dbConfig)
	spotHeroApp.Run(fmt.Sprintf(":%d",AppPort))
}