package base

import (
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func Log(s string, p ...interface{}) {
	log.Printf(s, p...)
}

func Debug(s string, p ...interface{}) {
	log.Printf(s, p...)
}

func LockGoRoutine() {
	for {
		time.Sleep(time.Minute)
	}
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func send404(w http.ResponseWriter, msg string) {
	w.WriteHeader(404)
	w.Write([]byte(msg))
}

func send500(w http.ResponseWriter) {
	w.WriteHeader(500)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Serve(roots []string) http.HandlerFunc {
	HFStatic := func(writer http.ResponseWriter, request *http.Request) {
		fname := strings.Split(request.URL.Path, "?")[0]
		if fname == "" || fname == "/" {
			fname = "/index.html"
		}
		for _, v := range roots {

			fqn := path.Join(v, fname)
			Debug("Getting: %s => %s", fname, fqn)
			if fileExists(fqn + ".gz") {
				Debug("FOUND: %s => %s", fname, fqn+".gz")
				ext := path.Ext(fqn)
				mimetype := mime.TypeByExtension(ext)
				writer.Header().Add("Content-Type", mimetype)
				writer.Header().Add("Content-Encoding", "gzip")
				bs, err := ioutil.ReadFile(fqn + ".gz")
				if err != nil {
					send500(writer)
					return
				}
				writer.Write(bs)
				return
			}
			if fileExists(fqn) {
				Debug("FOUND: %s => %s", fname, fqn)
				ext := path.Ext(fqn)
				mimetype := mime.TypeByExtension(ext)
				writer.Header().Add("Content-Type", mimetype)
				bs, err := ioutil.ReadFile(fqn)
				if err != nil {
					send500(writer)
					return
				}
				writer.Write(bs)
				return
			}
		}
		Debug("NOT FOUND: %s", fname)
		send404(writer, "Resource:"+fname+" not found")

	}

	return HFStatic
}

func GetEnv(s string, d string) string {
	ret := os.Getenv(s)
	if ret == "" {
		return d
	}
	return ret
}
