package gw

import (
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/forward"
	confighandler "github.com/digitalcircle-com-br/gw/lib/handlers/config"
	exechandler "github.com/digitalcircle-com-br/gw/lib/handlers/exec"
	execphandler "github.com/digitalcircle-com-br/gw/lib/handlers/execp"
	gorunhandler "github.com/digitalcircle-com-br/gw/lib/handlers/gorun"
	httphandler "github.com/digitalcircle-com-br/gw/lib/handlers/http"
	httpnocachehandler "github.com/digitalcircle-com-br/gw/lib/handlers/httpnocache"
	i18nhandler "github.com/digitalcircle-com-br/gw/lib/handlers/i18n"
	knockhandler "github.com/digitalcircle-com-br/gw/lib/handlers/knock"
	statichandler "github.com/digitalcircle-com-br/gw/lib/handlers/static"
	staticvhandler "github.com/digitalcircle-com-br/gw/lib/handlers/staticv"
	unixhandler "github.com/digitalcircle-com-br/gw/lib/handlers/unix"
	mw "github.com/digitalcircle-com-br/gw/lib/mw"
	"github.com/digitalcircle-com-br/gw/lib/util"
	"go.digitalcircle.com.br/lib/dynconfig"
	"go.digitalcircle.com.br/lib/env"
	"gopkg.in/yaml.v3"
)

type RPHandler interface {
	Name() string
	Init(first bool, c util.ConfigStruct, mux *http.ServeMux)
	CreateRP(pre string, s string) func(w http.ResponseWriter, r *http.Request)
}

var config util.ConfigStruct
var Handlers []RPHandler
var Mux http.ServeMux
var Serve func()

func CreateRPHandle(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {
	var ret func(w http.ResponseWriter, r *http.Request)
	for _, v := range Handlers {
		ret = v.CreateRP(pre, strurl)
		if ret != nil {
			return func(w http.ResponseWriter, r *http.Request) {
				ret(w, r)
			}
		}
	}
	return nil

}

func onConfig(isfirst bool, bs []byte) {
	//if config.Name != "" {
	//	dc.SetAppName(config.Name)
	//}
	err := yaml.Unmarshal(bs, &config)
	if err != nil {
		base.Log("Aborting config update Error processing config: %s", err.Error())
		return
	}
	Mux = http.ServeMux{}
	var pre string

	for k, v := range config.Env {
		base.Log("Setting ENV from Config: %s => %s", k, v)
		os.Setenv(k, v)
	}

	if config.Cors {
		base.Log("Using Cors: %s", env.GetD("CORS", "*"))
	}
	if config.Helmet {
		base.Log("Using helmet")
	}
	if config.Sts != "" {
		base.Log("Using STS: %s", config.Sts)
	}

	forward.Load(config.Forward)

	for k, v := range config.Routes {

		if strings.HasPrefix(k, "/") {
			pre = k
		} else {
			pre = strings.Join(strings.Split(k, "/")[1:], "/")
		}
		base.Log("Setting route: %s => %s", k, v)
		h := CreateRPHandle(pre, v)
		if h != nil {
			if config.Helmet {
				h = mw.Helmet(h)
			}
			if config.Cors {
				h = mw.CORS(h)
			}
			if config.Sts != "-" {
				h = mw.STS(config.Sts, h)
			}
			if config.Csp != "-" {
				h = mw.CSP(config.Csp, h)
			}
			if config.XFrame != "-" {
				h = mw.XFrame(config.XFrame, h)
			}
			Mux.HandleFunc(k, h)
		}

	}

	http.DefaultServeMux = &Mux
	config.Routes = make(map[string]string)
	for _, v := range Handlers {
		v.Init(isfirst, config, &Mux)
	}
	if isfirst {
		Serve()
	}
}

func RegisterHandler(handler RPHandler) {
	base.Log("Adding RPHandler: %s", handler.Name())
	Handlers = append(Handlers, handler)
}

var fname string
var saveMutex = sync.Mutex{}

func Init() {
	base.Log("GATEWAY - VER: %s", VER)

	RegisterHandler(&httphandler.HttpHandler{})
	RegisterHandler(&httpnocachehandler.HttpHandler{})
	RegisterHandler(&unixhandler.UnixHandler{})
	RegisterHandler(&knockhandler.KnockHandler{})
	RegisterHandler(&exechandler.ExecHandler{})
	RegisterHandler(&execphandler.ExecPHandler{})
	RegisterHandler(&statichandler.StaticHandler{})
	RegisterHandler(&staticvhandler.StaticVHandler{})
	RegisterHandler(&gorunhandler.GoRunHandler{})
	RegisterHandler(&confighandler.ConfigHandler{})
	RegisterHandler(&i18nhandler.I18NHandler{})

	wd, err := os.Getwd()

	util.ErrP(err)
	base.Log("Running from dir: %s", wd)
	fname = env.GetD("CONFIG", "./gw.yaml")
	wd = path.Dir(fname)
	base.Log("Loading config: %s", fname)
	dcopts := dynconfig.NewOpts()
	dcopts.Fname = fname
	dcopts.OnChange = onConfig
	dynconfig.Init(&config, dcopts)
	base.LockGoRoutine()
}
