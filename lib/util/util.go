package util

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type WebHook struct {
}

type ConfigStruct struct {
	//Direct Layer4 forward rules
	Forward map[string]string `yaml:forward`
	//Routes managed by this GW instance
	Routes map[string]string `yaml:"routes"`
	//Knock url for knock protected routes
	Knock string `yaml:"knock"`
	//Helmet is a time based penalty measure to mitigate DOS attackes
	Helmet bool `yaml:"helmet"`
	//Whether should use acme
	Acme bool `yaml:"acme"`
	//Addr to listen to
	SelfSigned bool `yaml:"self_signed"`
	//Addr to listen to
	Addr string `yaml:"addr"`
	//Where are the certs stored
	Certs string `yaml:"certs"`
	//Whether should go http or https
	Https bool `yaml:"https"`
	//Ena
	XFrame string `yaml:"xframe"`
	//In case no acme, this is the cert and key to be used
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
	//Env vars to be used elsewhere, complement the ones availablle, and overlap them.
	Env map[string]string `yaml:"env"`
	//CORS
	//https://developer.mozilla.org/en-US/docs/Glossary/CORS
	Cors     bool `yaml:"cors"`
	Insecure bool `yaml:"insecure"`
	//Strict-Transport-Security
	//STS https://developer.mozilla.org/pt-BR/docs/Web/HTTP/Headers/Strict-Transport-Security
	Sts string `yaml:"sts"`
	//Content-Security-Policy
	//https://developer.mozilla.org/en-US/docs/Web/HTTP/CSP
	Csp   string            `yaml:"csp"`
	Procs map[string]string `yaml:"procs"`
	//Cron  []dccron.Config   `yaml:"cron"`
	Name string `yaml:"name"`
	Wd   string `yaml:"wd"`
}

func GetAddr(r *http.Request) string {
	parts := strings.Split(r.RemoteAddr, ":")
	plen := len(parts)
	return strings.Join(parts[0:plen-1], ":")

}

func CreateRP(strurl string) *httputil.ReverseProxy {
	purl, err := url.Parse(strurl)
	if err != nil {
		panic(err)
	}
	proxy := httputil.NewSingleHostReverseProxy(purl)
	return proxy
}
func ErrP(err error) {
	if err != nil {
		panic(err)
	}
}
func Err(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Panic(err)
	}
}
