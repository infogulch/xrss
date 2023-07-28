package xrss

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	xtemplate "github.com/infogulch/caddy-xtemplate"
	"golang.org/x/exp/slices"
)

type XRss struct {
	Database struct {
		Driver  string `json:"driver,omitempty"`
		Connstr string `json:"connstr,omitempty"`
	} `json:"database,omitempty"`
	xtemplate *xtemplate.Templates
}

// Interface guards
var (
	_ caddy.Provisioner  = (*XRss)(nil)
	_ caddy.Validator    = (*XRss)(nil)
	_ caddy.CleanerUpper = (*XRss)(nil)

	_ caddyhttp.MiddlewareHandler = (*XRss)(nil)

	_ xtemplate.CustomFunctionsProvider = (*XRss)(nil)
)

// CustomTemplateFunctions implements xtemplate.CustomFunctionsProvider.
func (*XRss) CustomTemplateFunctions() template.FuncMap {
	return template.FuncMap{
		"parseRSS": funcParseRSS,
	}
}

// Validate ensures t has a valid configuration and implements caddy.Validator.
func (r *XRss) Validate() error {
	if r.Database.Driver != "" && slices.Index(sql.Drivers(), r.Database.Driver) == -1 {
		return fmt.Errorf("database driver '%s' does not exist", r.Database.Driver)
	}
	return nil
}

// Provision initialized the module and implements caddy.Provisioner.
func (r *XRss) Provision(ctx caddy.Context) error {
	if r.Database.Driver == "" {
		r.Database.Driver = "sqlite3"
		r.Database.Connstr = "file:rss.db?_journal=WAL&_vacuum=full&_foreign_keys=true&_synchronous=NORMAL"
	}
	r.xtemplate = &xtemplate.Templates{
		Root:        "templates",
		Database:    r.Database,
		FuncModules: []caddy.ModuleID{"http.handlers.xrss"},
	}
	return r.xtemplate.Provision(ctx)
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (r *XRss) ServeHTTP(w http.ResponseWriter, req *http.Request, next caddyhttp.Handler) error {
	return r.xtemplate.ServeHTTP(w, req, next)
}

// Cleanup  implements caddy.CleanerUpper.
func (r *XRss) Cleanup() error {
	return r.xtemplate.Cleanup()
}
