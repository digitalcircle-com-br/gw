package execp

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
)

type ExecPHandler struct {
}

func (k *ExecPHandler) Name() string {
	return "ExecPHandler"
}
func (k *ExecPHandler) Init(first bool, config util.ConfigStruct, mux *http.ServeMux) {

}
func (k *ExecPHandler) CreateRP(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(strurl, "execp://") {
		execroot := strings.Replace(strurl, "execp://", "", 1)
		return func(w http.ResponseWriter, r *http.Request) {
			base.Debug("Forwarding: %s => %s", r.URL.String(), strurl)
			r.Header.Add("X-FORWARD-FROM", r.RequestURI)
			r.RequestURI = strings.Replace(r.RequestURI, pre, "", 1)
			r.URL.Path = strings.Replace(r.URL.Path, pre, "", 1)
			exebname := strings.Split(r.URL.Path, "?")[0]
			exebname = strings.Replace(exebname, pre, "", 1)
			exename := path.Join(execroot, exebname)
			cmd := exec.Command(exename)
			cmd.Env = os.Environ()
			b := bytes.Buffer{}
			err := r.Write(&b)
			util.Err(err, w)
			pi, err := cmd.StdinPipe()
			util.Err(err, w)
			po, err := cmd.StdoutPipe()
			util.Err(err, w)

			err = cmd.Start()
			util.Err(err, w)

			pi.Write(b.Bytes())
			res, err := http.ReadResponse(bufio.NewReader(po), r)
			util.Err(err, w)
			for k, v := range res.Header {
				for _, v1 := range v {
					w.Header().Add(k, v1)
				}
			}
			w.WriteHeader(res.StatusCode)

			defer po.Close()

			io.Copy(w, res.Body)

		}
	}
	return nil
}
