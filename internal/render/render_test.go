package render

import (
	"github.com/Franlky01/bookingwebApp/internal/Models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td Models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session ")
	}
}
func TestRenderTemplate(t *testing.T) {
	pathToTemplate = "./../../templates"
	tc, err := CreateTestTemplateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter
	err = RenderTemplate(&ww, r, "home.page.gohtml", &Models.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}
	err = RenderTemplate(&ww, r, "non-existent.page.gohtml", &Models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}

}

//	result := AddDefaultData(&td,r)
//
//	if result == nil {
//		t.Error("failed")
//	}
//}
func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}
