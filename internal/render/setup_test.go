package render

import (
	"encoding/gob"
	"github.com/Franlky01/bookingwebApp/internal/Models"
	"github.com/Franlky01/bookingwebApp/internal/config"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	//what am i going to store in the session

	//this would store any value, its and interface
	gob.Register(Models.Reservation{})
	//change this to true when in Production

	testApp.InProduction = false

	//INFO LOGGER
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog
	//  these would print to terminal window ldate is a Date in a nice readable format
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()               // creating a new session
	session.Lifetime = 24 * time.Hour // how long should the session last
	session.Cookie.Persist = true     // cookie should persistent even  when browser is closed
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	//adding the session to the appconfig.
	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}

//these satisfies the interface

type myWriter struct{}

func (m *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (m *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
func (m *myWriter) WriteHeader(i int) {

}
