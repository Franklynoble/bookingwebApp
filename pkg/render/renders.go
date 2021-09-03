package render

import (
	"bytes"
	"fmt"
	"github.com/Franlky01/bookingwebApp/Models"
      "github.com/Franlky01/bookingwebApp/config"

	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var app *config.AppConfig

func AddDefaultData(td *Models.TemplateData) *Models.TemplateData {

	return td
}

//NewTemplate sets the config for the template package
func NewTemplate(a *config.AppConfig) {
	app = a // pass the parameter and assign it to pointer to a config.AppConfig
}

//Renders templates using html Template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *Models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache { //if UseCache is true, use the information from the template TemplateCache

		//get the template cache from the App config
		tc = app.TemplateCache
	} else { //else rebuild the templateCache
		tc, _ = CreateTemplateCache()
	}

	////creating template cache
	//tc, err := CreateTemplateCache()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//pull the template out of the Map, if ok it would return the map
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("err")
	}
	//write to hold bytes of buffer
	buf := new(bytes.Buffer) //creates new Buffer
	td = AddDefaultData(td)
	_ = t.Execute(buf, td)   //execute into the Buffer, also adding the TemplateData to the Buffer

	_, err := buf.WriteTo(w) // from the Buffer write to responseWriter
	if err != nil {
		fmt.Println("Error writing template to browser")
	}
}

//_, err := RenderTemplatetest(w)
//if err != nil {
//	fmt.Println("error getting template", err)
//	return
//}
//file to be changed
//
//	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
//	err = parsedTemplate.Execute(w, nil)
//	if err != nil {
//		fmt.Println("error parsing template:", err)
//		return
//	}
//}
//createTemplateCahed creates template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCach := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.gohtml")
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
		matches, err := filepath.Glob("./templates/*.layout.gohtml")
		if err != nil {
			return myCach, err
		}

		//test for the length
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.gohtml")
			if err != nil {

				return myCach, err
			}

		}
		myCach[name] = ts
	}

	return myCach, nil

}
