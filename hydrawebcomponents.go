package hydrawebcomponents

import (
	"net/http"
	"strings"
)

// Config configures a Web instance.
type Config struct {
	ProjectName string    // Full name, e.g. "HydraExperienceLibrary"
	BrandPrefix string    // Colored part, e.g. "Hydra"
	BrandSuffix string    // White part, e.g. "ExperienceLibrary"
	AdminToken  string    // Shared admin token for auth
	NavLinks    []NavLink // Navigation links shown when logged in
}

// NavLink is a navigation entry shown in the header.
type NavLink struct {
	Label string
	Path  string
}

// Brand holds the split project name for the header template.
type Brand struct {
	Prefix string
	Suffix string
}

// PageData is the base data wrapper that Render injects around project-specific data.
type PageData struct {
	Title    string
	Error    string
	LoggedIn bool
	Brand    Brand
	Nav      []NavLink
	Data     any // project-specific page data
}

// Web holds shared web infrastructure state.
type Web struct {
	adminToken string
	brand      Brand
	navLinks   []NavLink
	renderer   *Renderer
}

// New creates a new Web instance with the given configuration.
func New(cfg Config) *Web {
	return &Web{
		adminToken: cfg.AdminToken,
		brand: Brand{
			Prefix: cfg.BrandPrefix,
			Suffix: cfg.BrandSuffix,
		},
		navLinks: cfg.NavLinks,
	}
}

// IsAuthenticated checks if the request has a valid admin token
// via Bearer header or session cookie.
func (w *Web) IsAuthenticated(r *http.Request) bool {
	if auth := r.Header.Get("Authorization"); auth != "" {
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == w.adminToken {
			return true
		}
	}
	if cookie, err := r.Cookie("admin_session"); err == nil {
		if cookie.Value == w.adminToken {
			return true
		}
	}
	return false
}

// HandleLoginPage renders the shared login page.
func (w *Web) HandleLoginPage(wr http.ResponseWriter, r *http.Request) {
	if w.IsAuthenticated(r) {
		http.Redirect(wr, r, "/admin", http.StatusSeeOther)
		return
	}
	w.renderer.Render(wr, "login.html", nil, false, "")
}

// HandleLogin processes the login form submission.
func (w *Web) HandleLogin(wr http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if token != w.adminToken {
		wr.WriteHeader(http.StatusUnauthorized)
		w.renderer.Render(wr, "login.html", nil, false, "Invalid token")
		return
	}
	http.SetCookie(wr, &http.Cookie{
		Name:     "admin_session",
		Value:    w.adminToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 30,
	})
	http.Redirect(wr, r, "/admin", http.StatusSeeOther)
}

// HandleLogout clears the session cookie and redirects to login.
func (w *Web) HandleLogout(wr http.ResponseWriter, r *http.Request) {
	http.SetCookie(wr, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
	http.Redirect(wr, r, "/login", http.StatusSeeOther)
}

