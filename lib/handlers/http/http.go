package http

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
)

type HttpHandler struct {
}

func (k *HttpHandler) Name() string {
	return "HttpHanlder"
}
func (k *HttpHandler) Init(first bool, config util.ConfigStruct, mux *http.ServeMux) {

}
func (k *HttpHandler) CreateRP(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(strurl, "http://") || strings.HasPrefix(strurl, "https://") {
		rp := util.CreateRP(strurl)
		rp.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		ret := func(w http.ResponseWriter, r *http.Request) {

			base.Debug("Forwarding: %s => %s", r.URL.String(), strurl)
			r.Header.Add("X-FORWARD-FROM", r.RequestURI)
			r.RequestURI = strings.Replace(r.RequestURI, pre, "", 1)
			r.URL.Path = strings.Replace(r.URL.Path, pre, "", 1)

			rp.ServeHTTP(w, r)
		}

		return ret
	}
	return nil
}
