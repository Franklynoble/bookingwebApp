package main

import (
	"github.com/Franlky01/bookingwebApp/internal/config"
	"github.com/Franlky01/bookingwebApp/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	//mux := pat.New()//
	//mux.Get("/",http.HandlerFunc(handlers.Repo.Home))
	//mux.Get("/about",http.HandlerFunc(handlers.Repo.About))

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	mux.Get("/generals-quarters", handlers.Repo.Generals)

	mux.Get("/majors-suite", handlers.Repo.Majors)

	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)

	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/make-reservation", handlers.Repo.MakeReservations)
	mux.Post("/make-reservation", handlers.Repo.PostReservations)

	fileServer := http.FileServer(http.Dir("./static/"))
	//use mux to look for the path name
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
