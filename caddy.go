package xrss

import (
	"html/template"

	"github.com/caddyserver/caddy/v2"
)

func init() {
	caddy.RegisterModule(XRss{})
}

// CaddyModule returns the Caddy module information.
func (XRss) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "xtemplate.funcs.xrss",
		New: func() caddy.Module { return new(XRss) },
	}
}

type XRss struct{}

func (*XRss) Funcs() template.FuncMap {
	return template.FuncMap{
		"fetchFeed": funcFetchFeed,
	}
}

var (
	_ ExtraFuncs = (*XRss)(nil)
)

type ExtraFuncs interface {
	Funcs() template.FuncMap
}
