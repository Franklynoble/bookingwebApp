package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Franlky01/bookingwebApp/internal/Models"
	"github.com/Franlky01/bookingwebApp/internal/config"
	"github.com/Franlky01/bookingwebApp/internal/driver"
	"github.com/Franlky01/bookingwebApp/internal/forms"
	"github.com/Franlky01/bookingwebApp/internal/helpers"
	"github.com/Franlky01/bookingwebApp/internal/render"
	"github.com/Franlky01/bookingwebApp/internal/repository"
	"github.com/Franlky01/bookingwebApp/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"time"
)

//SWAPING COMPONENTWITHIN OUR APPLICATION WITH  MINIMAL
//CHANGES TO,USING REPOSITORY
//PATTERN is the effecient implementing this..
//Repo the repository used by the handlers
var Repo *Repository

//  Repository   type repository
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//NewRepo sets the Repository for the  handlers
//returns the instance of this type that gholds the application
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a, // create an instance of this type that holds the application
		DB:  dbrepo.NewpostgresRepo(db.SQL, a),
	}
}

//NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	//pull the remote address of the,
	remoteIP := r.RemoteAddr
	//put the remotadd to the session, using  a key to lookup the value and the value
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "home.page.gohtml", &Models.TemplateData{})
}

//about is the about handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	//stringMap["test"] = "hello, aGain"

	//add data Business Logic

	//pull the value out of the session
	//remoteIP := m.App.Session.GetString(r.Context(),"remote_ip")
	//    stringMap["remote_ip"] = remoteIP
	render.Template(w, r, "about.page.gohtml", &Models.TemplateData{
		StringMap: stringMap,
	})
}

//Generals renders the  generals room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.gohtml", &Models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "make-reservation.page.gohtml", &Models.TemplateData{})
}

//Majors renders the room Page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.gohtml", &Models.TemplateData{})
}

// Availability renders the  search availability Page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.gohtml", &Models.TemplateData{})
}

// PostAvailability renders the  search availability Page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	//every value gotten from req is string by default, and has to be converted to the type you want to use it for

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	if start == "" || end == "" {
		m.App.Session.Put(r.Context(), "error", "date can not be empty,select a valid date ")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	//    01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
	}
	//layout := "2006-01-02"
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	for _, i := range rooms {
		m.App.InfoLog.Println("ROOM:", i.ID, i.RoomName)
	}
	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability ")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
		//No Availability
		m.App.InfoLog.Println("No Availa..")

	}
	data := make(map[string]interface{})
	data["rooms"] = rooms
	// get start, and end date to be added to session
	res := Models.Reservation{

		StartDate: startDate,
		EndDate:   endDate,
	}
	//Add the res to Session
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.gohtml", &Models.TemplateData{
		Data: data,
	})

	w.Write([]byte(fmt.Sprintf(`start date is %s and end date is %s`, start, end)))
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	res := jsonResponse{
		Ok:      true,
		Message: "Available!",
	}
	out, err := json.MarshalIndent(res, "", "	")
	if err != nil {
		log.Println(err)
	}
	// create a header to tell the browser what kind of file you are expecting
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")

	w.Write(out)

}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.gohtml", &Models.TemplateData{})
}

func (m *Repository) MakeReservations(w http.ResponseWriter, r *http.Request) {
	var emptyReservation Models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Template(w, r, "make-reservation.page.gohtml", &Models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostReservations handles the posting of reservation form
func (m *Repository) PostReservations(w http.ResponseWriter, r *http.Request) {
	//render.render.Template(w, r, "make-reservation.page.gohtml", &Models.TemplateData{})

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	//date convertion, from string to date format, note every field received from the Browser is by default a string
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	//    01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	//layout := "2006-01-02"
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := Models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}
	form := forms.New(r.Form)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email")
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.gohtml", &Models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	NewReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	restriction := Models.RoomRestriction{

		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: NewReservationID,
		RestrictionID: 1,
	}
	err = m.DB.InsertRoomRestriction(restriction)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "reservation", reservation)

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(Models.Reservation)

	if !ok {
		log.Println("cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//Remove session after being used
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.gohtml", &Models.TemplateData{
		Data: data,
	})
}

func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	//convert string to int, the id would take exactly id from the routes
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//m.App.Session.Get(r.Context(),"reservation")

	res, ok := m.App.Session.Get(r.Context(), "reservation").(Models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
	}
	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
