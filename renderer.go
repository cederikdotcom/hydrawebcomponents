package hydrawebcomponents

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

// Renderer holds pre-parsed templates that combine the shared layout
// with project-specific page templates.
type Renderer struct {
	templates map[string]*template.Template
	web       *Web
}

// NewRenderer creates a Renderer by combining the shared layout template
// with project-specific page templates from projectFS.
//
// dir is the directory within projectFS (e.g. "templates").
// pages is the list of page template files to register (e.g. "admin.html").
// funcMap is an optional template.FuncMap for custom template functions.
func (w *Web) NewRenderer(projectFS fs.FS, dir string, pages []string, funcMap template.FuncMap) *Renderer {
	layoutBytes, err := fs.ReadFile(templates, "templates/layout.html")
	if err != nil {
		log.Fatalf("hydrawebcomponents: reading layout.html: %v", err)
	}
	loginBytes, err := fs.ReadFile(templates, "templates/login.html")
	if err != nil {
		log.Fatalf("hydrawebcomponents: reading login.html: %v", err)
	}

	r := &Renderer{
		templates: make(map[string]*template.Template),
		web:       w,
	}

	// Login uses shared layout + shared login template
	loginTmpl := template.New("layout.html")
	if funcMap != nil {
		loginTmpl = loginTmpl.Funcs(funcMap)
	}
	r.templates["login.html"] = template.Must(
		loginTmpl.Parse(string(layoutBytes) + "\n" + string(loginBytes)))

	// Project pages use shared layout + project content template
	for _, page := range pages {
		pageBytes, err := fs.ReadFile(projectFS, dir+"/"+page)
		if err != nil {
			log.Fatalf("hydrawebcomponents: reading %s/%s: %v", dir, page, err)
		}
		pageTmpl := template.New("layout.html")
		if funcMap != nil {
			pageTmpl = pageTmpl.Funcs(funcMap)
		}
		r.templates[page] = template.Must(
			pageTmpl.Parse(string(layoutBytes) + "\n" + string(pageBytes)))
	}

	// Store renderer reference so login/logout handlers can use it
	w.renderer = r

	return r
}

// Render executes the named template with project data wrapped in PageData.
func (r *Renderer) Render(w http.ResponseWriter, name string, data any, loggedIn bool, errMsg string) {
	t, ok := r.templates[name]
	if !ok {
		http.Error(w, "template not found: "+name, http.StatusInternalServerError)
		return
	}

	pd := PageData{
		Title:    r.web.brand.Prefix + r.web.brand.Suffix,
		LoggedIn: loggedIn,
		Brand:    r.web.brand,
		Nav:      r.web.navLinks,
		Error:    errMsg,
		Data:     data,
	}

	if err := t.ExecuteTemplate(w, "layout.html", pd); err != nil {
		log.Printf("hydrawebcomponents: template error (%s): %v", name, err)
	}
}
