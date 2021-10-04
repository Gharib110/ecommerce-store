package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type TemplateData struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRF            string
	Flash           string
	Warning         string
	Error           string
	IsAuthenticated int
	API             string
	CSSVersion      string
}

var functions = template.FuncMap{}

//go:embed templates
var templateFS embed.FS

// addDefaultData uses for adding default data into every TemplateData with http.Request
func (app *application) addDefaultData(td *TemplateData, r *http.Request) *TemplateData {
	return td
}

// renderTemplate uses for rendering the template that we need for presentation
func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request,
	page string, td *TemplateData, partials ...string) error {
	var t *template.Template
	var err error
	templateRender := fmt.Sprintf("templates/%s.page.tmpl", page)

	t, tmplOk := app.templateCache[templateRender]

	if app.config.env == "production" && tmplOk {
		t = app.templateCache[templateRender]
	} else {
		t, err = app.parseTemplate(page, templateRender, partials)
		if err != nil {
			app.errLogger.Println(err)
			return err
		}
	}

	if td == nil {
		td = &TemplateData{}
	}

	td = app.addDefaultData(td, r)
	err = t.Execute(w, td)
	if err != nil {
		return err
	}

	return nil
}

// parseTemplate uses for parsing page template with partials and other configurations
func (app *application) parseTemplate(page string,
	templateRender string, partials []string) (*template.Template, error) {
	var t *template.Template
	var err error

	// here I start building partials
	if len(partials) > 0 {
		for idx, partial := range partials {
			partials[idx] = fmt.Sprintf("templates/%s.partial.tmpl", partial)
		}
	}

	if len(partials) > 0 {
		t, err = template.New(fmt.Sprintf("%s.page.tmpl", page)).Funcs(functions).ParseFS(templateFS,
			"templates/base.layout.tmpl", strings.Join(partials, ","), templateRender)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.tmpl", page)).Funcs(functions).ParseFS(templateFS,
			"templates/base.layout.tmpl", templateRender)
	}

	if err != nil {
		app.errLogger.Println(err)
		return nil, err
	}

	app.templateCache[templateRender] = t

	return t, nil
}
