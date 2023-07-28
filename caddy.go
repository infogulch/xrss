package xrss

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	caddy.RegisterModule(XRss{})
}

// CaddyModule returns the Caddy module information.
func (XRss) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.xrss",
		New: func() caddy.Module { return new(XRss) },
	}
}

func init() {
	httpcaddyfile.RegisterHandlerDirective("xrss", parseCaddyfile)
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	rss := new(XRss)
	for h.Next() {
		for h.NextBlock(0) {
			switch h.Val() {
			case "database":
				for nesting := h.Nesting(); h.NextBlock(nesting); {
					switch h.Val() {
					case "driver":
						if !h.Args(&rss.Database.Driver) {
							return nil, h.ArgErr()
						}
					case "connstr":
						if !h.Args(&rss.Database.Connstr) {
							return nil, h.ArgErr()
						}
					}
				}
			}
		}
	}
	return rss, nil
}
