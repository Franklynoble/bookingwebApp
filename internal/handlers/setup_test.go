package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/Franlky01/bookingwebApp/internal/Models"
	"github.com/Franlky01/bookingwebApp/internal/config"
	"github.com/Franlky01/bookingwebApp/internal/render"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var pathToTemplate = "./../../templates"
var app config.AppConfig
var session *scs.SessionManager
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	//what am i going to store in the session
	//this would store any value, its and interface
	gob.Register(Models.Reservation{})
	//change this to true when in Production

	app.InProduction = false

	//INFO LOGGER
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	//  these would print to terminal window ldate is a Date in a nice readable format
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()               // creating a new session
	session.Lifetime = 24 * time.Hour // how long should the session last
	session.Cookie.Persist = true     // cookie should persistent even  when browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	//adding the session to the appconfig.
	app.Session = session
	tc, err := CreateTestTemplateCache()

	if err != nil {
		log.Fatal("can not create template cache")

	}

	app.TemplateCache = tc
	app.UseCache = true
	//created a new repository
	repo := NewRepo(&app)
	//pass it back to handlers
	NewHandlers(repo)

	render.NewTemplate(&app)

	//mux := pat.New()//
	//mux.Get("/",http.HandlerFunc(Repo.Home))
	//mux.Get("/about",http.HandlerFunc(Repo.About))

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)

	mux.Get("/generals-quarters", Repo.Generals)

	mux.Get("/majors-suite", Repo.Majors)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/make-reservation", Repo.MakeReservations)
	mux.Post("/make-reservation", Repo.PostReservations)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	fileServer := http.FileServer(http.Dir("./static/"))
	//use mux to look for the path name
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux

}

//NoSurf adds CSRF protection to all Post Request
func NoSurf(next http.Handler) http.Handler {
	crsfHanlder := nosurf.New(next)
	crsfHanlder.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/", // for perPage bases, the entire site
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return crsfHanlder
}

//loads and saves the session on every session
//tell the webapp to remember state using session
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

//createTemplateCahed creates template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCach := map[string]*template.Template{}
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplate))
	if err != nil {
		return myCach, err
	}

	//for every page found, the first page it would find using the index stating from 0 value
	// the pages here is the collection of  the file found in this variable

	for _, page := range pages {
		//return the full Path to the file
		name := filepath.Base(page) //here we get the base name

		fmt.Println("Page is currently on the page", page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCach, err

		}

		//if match it would be greater than zero
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplate))
		if err != nil {
			return myCach, err
		}

		//test for the length
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplate))
			if err != nil {
				return myCach, err
			}

		}
		myCach[name] = ts
	}

	return myCach, nil

}
