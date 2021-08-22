//+build dev

package httpz

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"selihc.com/glaive/buildinfo"
)

var devServerRunning bool
var devServerOnce sync.Once

func (r *Router) serveNotFound(w http.ResponseWriter, req *http.Request) {
	devServerOnce.Do(func() {
		fmt.Println("starting dev server")
		wd, err := os.Getwd()
		if err != nil {
			r.log.Error().Msg(err.Error())
		}

		fmt.Println(buildinfo.DevCommand)
		split := strings.Split(buildinfo.DevCommand, " ")
		cmd := exec.Command(split[0], split[1:]...)
		cmd.Dir = filepath.Join(wd, "/ui")
		cmd.Env = append(cmd.Env, os.Environ()...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Start()
		if err != nil {
			r.log.Error().Msg(err.Error())
		}

	})

	parsed, err := url.Parse(buildinfo.DevPort)
	if err != nil {
		r.log.Error().Msgf("%s, %w", buildinfo.DevPort, err)
		return
	}

	rp := httputil.NewSingleHostReverseProxy(parsed)
	rp.ServeHTTP(w, req)
}
