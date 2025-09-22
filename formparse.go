package formparse

import (
	"net/http"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

// Interface guards
var (
	_ caddy.Provisioner           = (*FormParse)(nil)
	_ caddyhttp.MiddlewareHandler = (*FormParse)(nil)
	_ caddyfile.Unmarshaler       = (*FormParse)(nil)
)

func init() {
	caddy.RegisterModule(FormParse{})
	httpcaddyfile.RegisterHandlerDirective("form_parse", parseCaddyfile)
}

// FormParse implements an HTTP handler that parses
// form body as placeholders.
type FormParse struct {
	FormKeys []string
	log      *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (FormParse) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.form_parse",
		New: func() caddy.Module { return new(FormParse) },
	}
}

// Provision implements caddy.Provisioner.
func (f *FormParse) Provision(ctx caddy.Context) error {
	f.log = ctx.Logger(f)

	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (f FormParse) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	repl := r.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)

	replacerFunc, err := newReplacerFunc(r, f.FormKeys)
	if err != nil {
		f.log.Debug("", zap.Error(err))
	}

	if err == nil {
		repl.Map(replacerFunc)
	}

	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (f *FormParse) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		f.FormKeys = d.RemainingArgs()
		if len(f.FormKeys) == 0 {
			return d.Errf("form_parse needs keys to be parsed")
		}
	}
	return nil
}

// parseCaddyfile unmarshals tokens from h into a new Middleware.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m FormParse
	err := m.UnmarshalCaddyfile(h.Dispenser)
	return m, err
}
