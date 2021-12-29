package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
	"gopkg.in/yaml.v3"
)

type ConfigHandler struct {
}

func (k *ConfigHandler) Name() string {
	return "ConfigHandler"
}
func (k *ConfigHandler) Init(first bool, config util.ConfigStruct, mux *http.ServeMux) {

}

func ErrH(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Fatal(err)
	}
}

func writeYAMLFileToJson(fname string, w http.ResponseWriter) error {
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
func (k *ConfigHandler) CreateRP(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(strurl, "config://") {

		dirname := strings.Replace(strurl, "config://", "", 1)
		deffname := path.Join(dirname, "default.yaml")
		return func(w http.ResponseWriter, r *http.Request) {

			orig := strings.Split(r.Host, ":")[0]
			fname := path.Join(dirname, orig+".yaml")

			if base.FileExists(fname) {
				log.Printf("Got origin %s, replying with: %s", orig, fname)
				err := writeYAMLFileToJson(fname, w)
				if err != nil {
					ErrH(err, w)
					return
				}

			} else if base.FileExists(deffname) {
				log.Printf("Got origin %s, replying with: %s", orig, deffname)
				err := writeYAMLFileToJson(deffname, w)
				if err != nil {
					ErrH(err, w)
					return
				}
			} else {
				log.Printf("Got origin %s, no file found: %s, %s", orig, fname, deffname)
				w.Write([]byte("{}"))
			}

		}
	}
	return nil
}
