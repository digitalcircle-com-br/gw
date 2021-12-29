package mw

import (
	"net/http"

	"github.com/digitalcircle-com-br/gw/lib/base"
)

func CSP(csp string, next http.HandlerFunc) http.HandlerFunc {

	if csp == "" || csp == "*" {
		csp = "default-src: 'self';"
	}
	base.Log("Setting CSP: %s", csp)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		(w).Header().Set("Content-Security-Policy:", csp)

		next.ServeHTTP(w, r)
	})
}
