package xrss

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	xtemplate "github.com/infogulch/caddy-xtemplate"
)

type XRss struct {
	xtemplate xtemplate.Templates
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
	if r.xtemplate.Database.Driver == "" {
		r.xtemplate.Database.Driver = "sqlite3"
		r.xtemplate.Database.Connstr = "file:rss.db?_journal=WAL&_synchronous=NORMAL&_foreign_keys=true&_vacuum=full"
	}
	r.xtemplate.TemplateRoot = "templates"
	r.xtemplate.ExtraFuncs = funcLibrary
	return r.xtemplate.Provision(ctx)
}

// Validate ensures t has a valid configuration and implements caddy.Validator.
func (r *XRss) Validate() error {
	return r.xtemplate.Validate()
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (r *XRss) ServeHTTP(w http.ResponseWriter, req *http.Request, next caddyhttp.Handler) error {
	return r.xtemplate.ServeHTTP(w, req, next)
}

// Cleanup  implements caddy.CleanerUpper.
func (r *XRss) Cleanup() error {
	return r.xtemplate.Cleanup()
}
