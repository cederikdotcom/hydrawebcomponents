package hydrawebcomponents

import (
	"html/template"
	"net/http"
)

// DocsSection represents a single documentation section (tab).
type DocsSection struct {
	Key   string
	Label string
}

// DocsSectionGroup is a labeled group of documentation sections.
type DocsSectionGroup struct {
	Label    string
	Sections []DocsSection
}

// DocsConfig configures the public documentation pages.
type DocsConfig struct {
	DefaultSection string                  // section key to redirect /docs to
	Tagline        string                  // shown on the landing page below the brand
	SectionGroups  []DocsSectionGroup       // tab groups for the docs viewer
	Cache          map[string]template.HTML // pre-rendered HTML keyed by section key
	EnableMermaid  bool                    // load Mermaid.js for diagram rendering
}

// docsPageData is the template data for the public docs viewer.
type docsPageData struct {
	SectionGroups []DocsSectionGroup
	ActiveSection string
	Content       template.HTML
}

// InitDocs configures the public docs/landing handlers. Call this after NewRenderer.
func (w *Web) InitDocs(cfg DocsConfig) {
	w.docsConfig = &cfg
}

// HandleLanding renders the landing page with Documentation and Admin Login buttons.
func (w *Web) HandleLanding(wr http.ResponseWriter, r *http.Request) {
	var tagline any
	if w.docsConfig != nil && w.docsConfig.Tagline != "" {
		tagline = w.docsConfig.Tagline
	}
	w.publicRenderer.Render(wr, "landing.html", tagline)
}

// HandleDocsRedirect redirects /docs to /docs/{default}.
func (w *Web) HandleDocsRedirect(wr http.ResponseWriter, r *http.Request) {
	section := "overview"
	if w.docsConfig != nil && w.docsConfig.DefaultSection != "" {
		section = w.docsConfig.DefaultSection
	}
	http.Redirect(wr, r, "/docs/"+section, http.StatusFound)
}

// HandleDocs renders a documentation section.
func (w *Web) HandleDocs(wr http.ResponseWriter, r *http.Request) {
	if w.docsConfig == nil {
		http.NotFound(wr, r)
		return
	}
	section := r.PathValue("section")
	content, ok := w.docsConfig.Cache[section]
	if !ok {
		http.NotFound(wr, r)
		return
	}
	w.publicRenderer.Render(wr, "public-docs.html", docsPageData{
		SectionGroups: w.docsConfig.SectionGroups,
		ActiveSection: section,
		Content:       content,
	})
}
