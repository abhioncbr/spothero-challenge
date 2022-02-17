package app

import (
	"log"
	"net/http"
	"spotHero/app/handler"
	"spotHero/app/model"
	"spotHero/config"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(dbConfig *config.DBConfig) {

	db, err := gorm.Open(dbConfig.GormDialect, &dbConfig.GormConfig)
	if err != nil {
		log.Fatal("Could not connect database", err)
	}

	_ = model.DBMigrate(db)
	dataLoadErr := model.LoadRatesOnStart("rates.json" , db)
	if dataLoadErr != nil {
		log.Fatal("Could not load rate list in to DB, error: ", dataLoadErr)
	}

	a.DB = db
	a.Router = mux.NewRouter()
	a.setRouters()
}

// setRouters sets the all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/rates", a.handleRequest(handler.GetAllRates))
	a.Put("/rates", a.handleRequest(handler.PutRate))
	a.Get("/price", a.handleRequest(handler.GetPrice))
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}