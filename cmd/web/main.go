package main

import (
	"fmt"
	"github.com/Franlky01/bookingwebApp/config"

	"github.com/Franlky01/bookingwebApp/pkg/handlers"
	"github.com/Franlky01/bookingwebApp/pkg/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	//change this to true when in Production
	app.InProduction = false
	session = scs.New()               // creating a new session
	session.Lifetime = 24 * time.Hour // how long should the session last
	session.Cookie.Persist = true     // cookie should persistent even  when browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	//adding the session to the appconfig.
	app.Session = session
	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("can not create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false
	//created a new repository
	repo := handlers.NewRepo(&app)
	//pass it back to handlers
	handlers.NewHandlers(repo)

	render.NewTemplate(&app)

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
