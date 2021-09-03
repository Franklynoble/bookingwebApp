package handlers

import (
	"github.com/Franlky01/bookingwebApp/Models"
	"github.com/Franlky01/bookingwebApp/config"
	"github.com/Franlky01/bookingwebApp/pkg/render"
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
	 m.App.Session.Put(r.Context(),"remote_ip", remoteIP)


	render.RenderTemplate(w, "home.page.gohtml",&Models.TemplateData{})
}

//about is the about handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)


	remoteIP := m.App.Session.GetString(r.Context(),"remote_ip")
	   stringMap["remote_ip"] = remoteIP
	//stringMap["test"] = "hello, aGain"

	//add data Business Logic

	//pull the value out of the session
	//remoteIP := m.App.Session.GetString(r.Context(),"remote_ip")
 //    stringMap["remote_ip"] = remoteIP
	 render.RenderTemplate(w, "about.page.gohtml",  &Models.TemplateData{
		StringMap: stringMap,
	})

}
