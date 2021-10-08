package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Franlky01/bookingwebApp/internal/config"
	"github.com/Franlky01/bookingwebApp/internal/driver"
	"github.com/Franlky01/bookingwebApp/internal/forms"
	"github.com/Franlky01/bookingwebApp/internal/helpers"
	"github.com/Franlky01/bookingwebApp/internal/models"
	"github.com/Franlky01/bookingwebApp/internal/render"
	"github.com/Franlky01/bookingwebApp/internal/repository"
	"github.com/Franlky01/bookingwebApp/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingsRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "about.page.gohtml", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find room!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//build  the reservation  Model up
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	// generate a form and  start testing
	form := forms.New(r.PostForm)

	// test all the fields from the user Form
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3) // first name is atleast three character
	form.IsEmail("email")           // make sure email is in form of email

	if !form.Valid() { // if any field in form fails , return error and retain the Form also highlight  the error that they Have
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "my own error message", http.StatusSeeOther)
		render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// if it gets passed that form input start,  then insert the details into the Database
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return

	}

	// send notifications - First to Guest
	htmlMessage := fmt.Sprintf(`<strong>Reservation Confirmation</strong><br>

Dear %s, <br>
This is to confirm your reservation from %s to %s,`,
		reservation.FirstName, reservation.StartDate.Format("2006-01.02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:      reservation.Email,
		From:    "me@here.com",
		Subject: "Reservation Confirmation",
		Content: htmlMessage,
	}
	m.App.MailChan <- msg // pass the Message for mail channel

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.gohtml", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.gohtml", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	// pass form populates the r.form for post request, it passes the raw request from the URl and Updates the Form
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(rooms) == 0 {
		// no availability
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	// need to parse request body
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	roomID, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		// got a database error, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}

		// removed the error check, since we handle all aspect of the Json right here
		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	// I removed the error check, since we handle all aspects of
	// the json right here
	out, _ := json.MarshalIndent(resp, "", "     ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// Contact renders the search availability page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.gohtml", &models.TemplateData{})
}

// ReservationSummary displays the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// used to have next 6 lines
	//roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	//if err != nil {
	//	log.Println(err)
	//	m.App.Session.Put(r.Context(), "error", "missing url parameter")
	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	//	return
	//}

	// changed to this, so we can test it more easily
	// split the URL up by /, and grab the 3rd element
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// BookRoom takes URL parameters, builds a sessional variable, and takes user to make res screen
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	if err != nil {
		log.Println(err, "the id not passed")
	}
	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't get room from db!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.gohtml", &models.TemplateData{
		Form: forms.New(nil)})
}

// PostShowLogin handles logging in the user
func (m Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	//doing login and log out it is good practice to renew the code
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	// using post form, you use this,
	// get parameter from the form
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	form := forms.New(r.PostForm)
	//authenticating fields you use this fields
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		//TODO -take user back to page
		render.Template(w, r, "login.page.gohtml", &models.TemplateData{
			Form: form,
		})
		return
	}
	//now authenticate for  the user
	id, _, err := m.DB.Authenticate(email, password)

	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	//if successfully loged in pu the session back
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//  Logout logs a user out
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	//to logout, you must make sure the token is moved from the session, to do that you destroy the session,and renew it
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

//AdminDashboard shows the admin dashboard
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.gohtml", &models.TemplateData{})
}

//AdminNewReservations shows  all new reservations in Admin Tool
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllNewReservations()

	if err != nil {
		log.Println(err)
	}
	data := make(map[string]interface{})

	data["reservations"] = reservations

	render.Template(w, r, "admin-new-reservations.page.gohtml", &models.TemplateData{
		Data: data,
	})

}
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllReservations()

	if err != nil {
		log.Println(err)
	}
	data := make(map[string]interface{})

	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.gohtml", &models.TemplateData{
		Data: data,
	})
}

//AdminReservationsCalender displays the reservations Calender
func (m *Repository) AdminReservationsCalender(w http.ResponseWriter, r *http.Request) {
	//assume that there is  no month specified

	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))

		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	}
	//dataMap for date textual format

	data := make(map[string]interface{})

	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)
	//two digits Month Format
	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)

	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear

	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	// format the date passed  ...

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get the First and last days of the month

	currenttYear, currentMonth, _ := now.Date()

	currentLocation := now.Location()
	firstOfMonth := time.Date(currenttYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)

	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data["rooms"] = rooms

	for _, x := range rooms {
		//create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}
		// get  all restrictions for the current room
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		for _, y := range restrictions {
			if y.ReservationID > 0 {
				//it's a  reservations
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
			} else {
				//it's a block
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID
			}
		}
		//this is registered in the glob registere
		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	render.Template(w, r, "admin-reservations-calender.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

//AdminShowReservations shows  all new reservations in Admin Tool
func (m *Repository) AdminShowReservations(w http.ResponseWriter, r *http.Request) {
	//get the URI variable at index position and  explode it i.e, select the vairaiable based on the index,
	//index position count starts at 1
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4]) // get it at index position 4

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3] // get element from URI at index position 3

	stringMap := make(map[string]string)
	stringMap["src"] = src

	//get reservation from  the database
	reservation, _ := m.DB.GetReservationByID(id) // if id not gotten, it would show  the  nil pointer address
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservation"] = reservation // pass reservation id from db to the map

	render.Template(w, r, "admin-reservations-show.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

func (m *Repository) AdminPostReservation(w http.ResponseWriter, r *http.Request) {
	// pass form populates the r.form for post request, it passes the raw request from the URl and Updates the Form
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		helpers.ServerError(w, err)
		return
	}
	//get the URI variable at index position and  explode it i.e, select the vairaiable based on the index,
	//index position count starts at 1
	exploded := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(exploded[4]) // get it at index position 4

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3] // get element from URI at index position 3

	stringMap := make(map[string]string)
	stringMap["src"] = src
	res, err := m.DB.GetReservationByID(id)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = m.DB.UpdatReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-#{src}"), http.StatusSeeOther)
}

// AdminProcessReservation marks a reservation calendar
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	fmt.Println(id, src)
	_ = m.DB.UpdateProcessedForReservation(id, 1)
	m.App.Session.Put(r.Context(), "flash", "Reservation Marked as Processd")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}

//AdminDeleteReservation deletes a Reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	fmt.Println(id, src)
	_ = m.DB.DeleteReservation(id)
	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}
func (m *Repository) AdminPostReservationsCalender(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
	}

	//get date from the URL
	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	//process blocks

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)

	for _, x := range rooms {
		//			get the block map from the session, Loop through entire map, if we have an entry in the map
		//	that does not exist in our posted, and if the restriction id > 0, then it is a block we need to
		//remove

		// get the session context
		// cast the block_map into map,strint int
		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range curMap {
			// ok , will be false if the value is not in the map
			if val, ok := curMap[name]; ok {
				// only pay attention to value > 0, and that are not in the form post
				// the rest are just placeholders for days without block
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {
						// delete the restriction by id
						//	log.Println("would delete block", value)
						err := m.DB.DeleteBlockByID(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}
	// now handle new blocks
	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_") // splits back  a sclice of string
			roomID, _ := strconv.Atoi(exploded[2])

			//insert  a new Block
			//	log.Println("Would insert a block id ", roomID , "for date",exploded[3])
			t, _ := time.Parse("2006-01-2", exploded[3])
			err := m.DB.InsertBlockForRoom(roomID, t)
			if err != nil {
				log.Println(err)
			}

		}

		//log.Println("form has Name:",name)
	}
	m.App.Session.Put(r.Context(), "flash", "changes saved")
	http.Redirect(w, r, fmt.Sprintf("admin/reservations-calender?y=%d&m=%d", year, month), http.StatusSeeOther)
}
