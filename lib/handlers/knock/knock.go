package knock

import (
	"net/http"
	"strings"
	"time"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
)

var knock map[string]int64

func regKnock(w http.ResponseWriter, r *http.Request) {
	addr := util.GetAddr(r)
	base.Debug("Knock from %s", addr)
	knock[addr] = time.Now().Unix()
	w.WriteHeader(200)
	w.Write([]byte("Knock-Knock! - " + addr))
}
func regUnKnock(w http.ResponseWriter, r *http.Request) {
	addr := util.GetAddr(r)
	delete(knock, addr)
	w.WriteHeader(200)
	w.Write([]byte("Bye! - " + addr))
}
func checkKnock(r *http.Request) bool {
	addr := util.GetAddr(r)
	_, ok := knock[addr]
	return ok
}
func cleanKnock() {
	for {
		now := time.Now().Unix()

		keys := make([]string, 0)

		for k, v := range knock {
			if now-v > 60*60*6 {
				base.Debug("Cleaning up knock from %s", k)
				keys = append(keys, k)
			}
		}

		for _, v := range keys {
			delete(knock, v)
		}

		//delete(knock, k)
		time.Sleep(5 * time.Second)
	}
}

type KnockHandler struct {
}

func (k *KnockHandler) Name() string {
	return "KnockHandler"
}
func (k *KnockHandler) Init(first bool, config util.ConfigStruct, mux *http.ServeMux) {
	if first {
		knock = make(map[string]int64)
	}

	if config.Knock != "" {
		base.Log("Enabling Knock at: %s", config.Knock)
		mux.HandleFunc(config.Knock+"/hello", regKnock)
		mux.HandleFunc(config.Knock+"/bye", regUnKnock)
		go cleanKnock()
	}

}
func (k *KnockHandler) CreateRP(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(strurl, "knock://") {
		strurl = strings.Replace(strurl, "knock://", "", 1)
		rp := util.CreateRP(strurl)

		return func(w http.ResponseWriter, r *http.Request) {

			checkKnock(r)
			if !checkKnock(r) {
				time.Sleep(time.Second * 10)
				base.Debug("Negating Knock from: %s ", r.RemoteAddr)
				w.WriteHeader(403)
				w.Write([]byte("Who's there?"))
				return
			}

			base.Debug("Forwarding: %s => %s", r.URL.String(), strurl)
			r.Header.Add("X-FORWARD-FROM", r.RequestURI)
			r.RequestURI = strings.Replace(r.RequestURI, pre, "", 1)
			r.URL.Path = strings.Replace(r.URL.Path, pre, "", 1)
			rp.ServeHTTP(w, r)
		}
	}
	return nil
}
