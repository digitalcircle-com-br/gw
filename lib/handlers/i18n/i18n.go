package i18n

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
	"gopkg.in/yaml.v3"
)

type I18NHandler struct {
}

func (k *I18NHandler) Name() string {
	return "I18NHandler"
}
func (k *I18NHandler) Init(first bool, config util.ConfigStruct, mux *http.ServeMux) {

}

func ErrH(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatal(err)
	}
}

func writeYAMLFileToJson(fname string, w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "application/json")
	c := make(map[string]interface{})
	bs, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bs, &c)
	if err != nil {
		return err
	}
	jsonbs, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Write(jsonbs)
	return nil
}
func (k *I18NHandler) CreateRP(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(strurl, "i18n://") {

		dirname := strings.Replace(strurl, "i18n://", "", 1)
		deffname := path.Join(dirname, "default.yaml")
		return func(w http.ResponseWriter, r *http.Request) {
			fname := deffname
			l := r.URL.Query().Get("l")
			if l == "" {
				l = strings.ToLower(strings.Split(r.Header.Get("Accept-Language"), ",")[0])
			}
			if l == "" {
				l = strings.Split(os.Getenv("LANG"), ".")[0]
			}

			if l != "" {
				fname = l + ".yaml"
			}

			if base.FileExists(path.Join(dirname, fname)) {
				base.Debug("Got Lang %s, replying with: %s", l, fname)
				err := writeYAMLFileToJson(path.Join(dirname, fname), w)
				if err != nil {
					ErrH(err, w)
					return
				}

			} else {
				base.Debug("Got origin %s, no file found: %s, %s", l, fname, deffname)

				w.Write([]byte("{}"))
			}

		}
	}
	return nil
}
