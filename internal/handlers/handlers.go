package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/Franlky01/bookingwebApp/internal/Models"
	"github.com/Franlky01/bookingwebApp/internal/config"
	"github.com/Franlky01/bookingwebApp/internal/render"
	"log"
	"net/http"
)

//SWAPING COMPONENTWITHIN OUR APPLICATION WITH  MINIMAL
//CHANGES TO,USING REPOSITORY
//PATTERN is the effecient implementing this..
//Repo the repository used by the handlers
var Repo *Repository

//  Repository   type repository
type Repository struct {
	App *config.AppConfig
}

//NewRepo sets the Repository for the  handlers
//returns the instance of this type that gholds the application
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a, // create an instance of this type that holds the application
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

	render.RenderTemplate(w, r, "home.page.gohtml", &Models.TemplateData{})
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
	render.RenderTemplate(w, r, "about.page.gohtml", &Models.TemplateData{
		StringMap: stringMap,
	})
}

//Generals renders the  generals room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.gohtml", &Models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &Models.TemplateData{})
}

//Majors renders the room Page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.gohtml", &Models.TemplateData{})
}

// Availability renders the  search availability Page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &Models.TemplateData{})
}

// PostAvailability renders the  search availability Page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	//every value gotten from req is string by default, and has to be converted to the type you want to use it for

	start := r.Form.Get("start")
	end := r.Form.Get("end")

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
	render.RenderTemplate(w, r, "contact.page.gohtml", &Models.TemplateData{})
}

func (m *Repository) MakeReservations(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &Models.TemplateData{})
}
