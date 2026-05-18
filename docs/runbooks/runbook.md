# hydrawebcomponents runbook

Shared Go library providing HTML templates, auth middleware, login flow, docs rendering, and layout components for Hydra web services.

## Overview

This is a Go library — it has no runtime binary or systemd service. It is consumed by other Hydra services (hydravenues, hydraneck, hydraexperiencelibrary, etc.) via `go get`.

## Releasing a new version

```bash
git tag v<X.Y.Z>
git push origin v<X.Y.Z>
```

The Go module proxy picks up the tag automatically. Consumers update by running:

```bash
go get github.com/cederikdotcom/hydrawebcomponents@v<X.Y.Z>
go mod tidy
```

## Key components

| File | Purpose |
|------|---------|
| `hydrawebcomponents.go` | `Web` struct, `HandleLoginPage`, `HandleLogin`, `HandleLogout` |
| `middleware.go` | `RequireAuth` (Bearer), `RequireWebAuth` (cookie + ?next= redirect), `LogRequest` |
| `renderer.go` | `Renderer` — wraps project templates in the shared layout |
| `templates/login.html` | Shared login page — passes `next` hidden field through the form |
| `templates/layout.html` | Shared admin layout (nav, brand) |
| `templates/public-layout.html` | Unauthenticated public layout |
| `docs.go` | `DocsConfig`, `HandleDocs`, `HandleDocsRedirect` |
| `markdown.go` | `RenderMarkdown` — goldmark-based Markdown → HTML |

## Authentication flow

1. `RequireWebAuth` checks the session cookie via `hydraauth`.
2. On failure it redirects to `/login?next=<original-url>`.
3. `HandleLoginPage` renders login form with a hidden `next` field.
4. `HandleLogin` validates the token, sets the cookie, then redirects to `next` (must be a relative path — host-prefixed values are rejected to prevent open redirect).
5. If `next` is empty, redirect goes to `/admin`.

## Integrating into a new service

```go
import hwc "github.com/cederikdotcom/hydrawebcomponents"

w := hwc.New(hwc.Config{
    BrandPrefix: "Hydra",
    BrandSuffix: "MyService",
    Auth:        myAuthInstance,
    NavLinks:    []hwc.NavLink{{Label: "Dashboard", Path: "/admin"}},
})
w.NewRenderer(web.Templates, "templates", []string{"admin.html"}, nil)

mux.HandleFunc("GET /login",  w.HandleLoginPage)
mux.HandleFunc("POST /login", w.HandleLogin)
mux.HandleFunc("POST /logout", w.HandleLogout)
mux.HandleFunc("GET /admin", w.RequireWebAuth(myHandler))
```

Project templates use the `{{define "content"}}` block. Access page data via `{{.Data}}` and error message via `{{.Error}}`.

## Troubleshooting

**Login redirects to /admin instead of the original page**
The `?next=` param was not preserved. Check that:
- The service uses `w.RequireWebAuth` (from this library), not `auth.RequireWebAuth` directly.
- The login form template has the hidden `next` input (present in `templates/login.html` since v0.6.0).

**Template not found error**
The page template was not passed in the `pages` slice to `NewRenderer`. Either add it there or register it separately.

**Login page shows blank after failed token**
`HandleLogin` re-renders `login.html` with the `next` value from the form. If the form does not include the hidden `next` field, the redirect after a retry goes to `/admin` instead of the original page.
