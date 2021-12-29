package unix

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
)

type UnixHandler struct {
}

func (k *UnixHandler) Name() string {
	return "UnixHandler"
}
func (k *UnixHandler) Init(first bool, config util.ConfigStruct, mux *http.ServeMux) {

}
func (k *UnixHandler) CreateRP(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(strurl, "unix://") {
		sock := strings.Replace(strurl, "unix://", "", 1)
		return func(w http.ResponseWriter, r *http.Request) {
			base.Debug("Forwarding: %s => %s", r.URL.String(), strurl)
			r.Header.Add("X-FORWARD-FROM", r.RequestURI)
			r.RequestURI = strings.Replace(r.RequestURI, pre, "", 1)
			r.URL.Path = strings.Replace(r.URL.Path, pre, "", 1)
			b := bytes.Buffer{}
			r.Write(&b)
			conn, err := net.Dial("unix", sock)
			if err != nil {
				base.Log(err.Error())
				http.Error(w, err.Error(), 500)
				return
			}
			defer conn.Close()

			go func() {
				conn.Write(b.Bytes())
			}()

			res, err := http.ReadResponse(bufio.NewReader(conn), r)

			if err != nil {
				base.Log(err.Error())
				http.Error(w, err.Error(), 500)
				return
			}
			for k, vs := range res.Header {
				for _, v := range vs {
					w.Header().Add(k, v)
				}

			}
			w.WriteHeader(res.StatusCode)
			io.Copy(w, res.Body)
			return
		}
	}
	return nil

}
