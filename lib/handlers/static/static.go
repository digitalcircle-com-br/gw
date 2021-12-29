package static

import (
	"net/http"
	"strings"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
)

type StaticHandler struct {
}

func (k *StaticHandler) Name() string {
	return "StaticHandler"
}
func (k *StaticHandler) Init(first bool, config util.ConfigStruct, mux *http.ServeMux) {

}
func (k *StaticHandler) CreateRP(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(strurl, "static://") {
		pathsstr := strings.Replace(strurl, "static://", "", 1)
		pathraw := strings.Split(pathsstr, ",")
		paths := make([]string, 0)
		for _, v := range pathraw {
			paths = append(paths, strings.TrimSpace(v))
		}

		statichandler := base.Serve(paths)
		ret := func(w http.ResponseWriter, r *http.Request) {
			base.Debug("Forwarding: %s => %s", r.URL.String(), strurl)
			r.Header.Add("X-FORWARD-FROM", r.RequestURI)
			r.RequestURI = strings.Replace(r.RequestURI, pre, "", 1)
			r.URL.Path = strings.Replace(r.URL.Path, pre, "", 1)
			statichandler(w, r)
		}

		return ret
	}
	return nil
}
