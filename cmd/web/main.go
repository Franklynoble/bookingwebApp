package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Franlky01/bookingwebApp/internal/Models"
	"github.com/Franlky01/bookingwebApp/internal/config"
	"github.com/Franlky01/bookingwebApp/internal/driver"
	"github.com/Franlky01/bookingwebApp/internal/handlers"
	"github.com/Franlky01/bookingwebApp/internal/helpers"
	"github.com/Franlky01/bookingwebApp/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	gob.Register(Models.Reservation{})

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()
	//give a function receiver to access the variable
	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)
	fmt.Println(fmt.Sprintf("starting application at %s", portNumber))
	//http.ListenAndServe(portNumber, nil)

	serve := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = serve.ListenAndServe()
	log.Fatal(err)

}
func run() (*driver.DB, error) {

	//change this to true when in Production
	app.InProduction = false
	session = scs.New()               // creating a new session
	session.Lifetime = 24 * time.Hour // how long should the session last
	session.Cookie.Persist = true     // cookie should persistent even  when browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	//adding the session to the appconfig.
	app.Session = session
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	log.Println("Connecting to Database ....")

	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=postgres user=postgres password=frank")

	if err != nil {
		log.Fatal("Can not connect to database Dying...")
	}
	log.Println("Connected to Database!")

	tc, err := render.CreateTestTemplateCache()

	if err != nil {
		log.Fatal("can not create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false
	//created a new repository
	repo := handlers.NewRepo(&app, db)
	//pass it back to handlers
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
