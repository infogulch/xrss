package xrss

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	xtemplate "github.com/infogulch/caddy-xtemplate"
)

type XRss struct {
	xtemplate.Templates
}

// Interface guards
var (
	_ caddy.Provisioner  = (*XRss)(nil)
	_ caddy.Validator    = (*XRss)(nil)
	_ caddy.CleanerUpper = (*XRss)(nil)

	_ caddyhttp.MiddlewareHandler = (*XRss)(nil)
)

// Provision initialized the module and implements caddy.Provisioner.
func (r *XRss) Provision(ctx caddy.Context) error {
	if r.Templates.Database.Driver == "" {
		r.Templates.Database.Driver = "sqlite3"
		r.Templates.Database.Connstr = "file:rss.sqlite?_journal=WAL&_synchronous=NORMAL&_foreign_keys=true&_vacuum=full"
	}
	r.Templates.TemplateRoot = "templates"
	r.Templates.ContextRoot = "schema"
	r.Templates.Delimiters = []string{"{{", "}}"}
	r.Templates.ExtraFuncs = funcLibrary
	return r.Templates.Provision(ctx)
}

// Validate ensures t has a valid configuration and implements caddy.Validator.
func (r *XRss) Validate() error {
	return r.Templates.Validate()
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (r *XRss) ServeHTTP(w http.ResponseWriter, req *http.Request, next caddyhttp.Handler) error {
	return r.Templates.ServeHTTP(w, req, next)
}

// Cleanup  implements caddy.CleanerUpper.
func (r *XRss) Cleanup() error {
	return r.Templates.Cleanup()
}
