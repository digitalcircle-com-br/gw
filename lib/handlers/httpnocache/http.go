package httpnocache

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
)

type HttpHandler struct {
}

func (k *HttpHandler) Name() string {
	return "HttpNoCacheHanlder"
}
func (k *HttpHandler) Init(first bool, config util.ConfigStruct, mux *http.ServeMux) {

}
func (k *HttpHandler) CreateRP(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(strurl, "httpnc://") || strings.HasPrefix(strurl, "httpncs://") {
		strurl = strings.Replace(strurl, "httpnc", "http", 1)
		rp := util.CreateRP(strurl)
		rp.FlushInterval = -1
		rp.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		ret := func(w http.ResponseWriter, r *http.Request) {

			f, ok := w.(http.Flusher)
			if ok {
				go func() {
					defer func() {
						r := recover()
						if r != nil {
							log.Printf("ERR at Flush: %v", r)
						}
					}()
					for {
						f.Flush()
						time.Sleep(time.Millisecond * 10)
					}

				}()
			}

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
