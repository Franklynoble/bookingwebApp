package config

import (
	"github.com/Franlky01/bookingwebApp/internal/models"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

//AppConfig holds application config, datas here would be exposed to all parts of the Application
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger //
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
