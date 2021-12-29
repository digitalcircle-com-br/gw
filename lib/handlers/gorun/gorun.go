package gorun

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
)

type GoRunHandler struct {
}

func (k *GoRunHandler) Name() string {
	return "GoRunHandler"
}
func (k *GoRunHandler) Init(first bool, config util.ConfigStruct, mux *http.ServeMux) {

}
func (k *GoRunHandler) CreateRP(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(strurl, "gorun://") {
		exename := strings.Replace(strurl, "gorun://", "", 1)
		return func(w http.ResponseWriter, r *http.Request) {
			base.Debug("Forwarding: %s => %s", r.URL.String(), strurl)
			r.Header.Add("X-FORWARD-FROM", r.RequestURI)
			r.RequestURI = strings.Replace(r.RequestURI, pre, "", 1)
			r.URL.Path = strings.Replace(r.URL.Path, pre, "", 1)
			if !strings.HasPrefix(r.URL.Path, "/") {
				r.URL.Path = "/" + r.URL.Path
			}
			if !strings.HasPrefix(r.RequestURI, "/") {
				r.RequestURI = "/" + r.RequestURI
			}
			shell := base.GetEnv("SHELL", "/bin/sh")
			cmd := exec.Command(shell, "-c", "go run "+exename)
			cmd.Env = os.Environ()
			cmd.Stderr = os.Stderr
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
			if err != nil {
				log.Printf("%s", err.Error())
				http.Error(w, err.Error(), 500)
				return
			}
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
