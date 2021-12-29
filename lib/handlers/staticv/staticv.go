package staticv

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/digitalcircle-com-br/gw/lib/base"
	"github.com/digitalcircle-com-br/gw/lib/util"
	"github.com/fsnotify/fsnotify"
)

type StaticVHandler struct {
}

func (k *StaticVHandler) Name() string {
	return "StaticVHandler"
}
func (k *StaticVHandler) Init(first bool, config util.ConfigStruct, mux *http.ServeMux) {

}
func (k *StaticVHandler) CreateRP(pre string, strurl string) func(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(strurl, "staticv://") {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}

		watchDir := func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// since fsnotify can watch all the files in a directory, watchers only need
			// to be added to each nested directory
			if fi != nil {
				if fi.Mode().IsDir() {
					base.Log("Watching: %s", path)
					return watcher.Add(path)
				}
			}

			return nil
		}

		go func() {
			for {
				select {
				// watch for events
				case event := <-watcher.Events:
					fmt.Printf("EVENT! %#v\n", event)

					// watch for errors
				case err := <-watcher.Errors:
					fmt.Println("ERROR", err)
				}
			}
		}()

		pathsstr := strings.Replace(strurl, "staticv://", "", 1)
		pathraw := strings.Split(pathsstr, ",")
		paths := make([]string, 0)
		for _, v := range pathraw {
			v = strings.TrimSpace(v)
			paths = append(paths, v)
			if strings.HasPrefix(v, ".") {
				v = path.Join(wd, v)
			}
			watcher.Add(v)
			err = filepath.Walk(v, watchDir)
			if err != nil {
				log.Printf("Error at staticv:%s", err.Error())
			}

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
