package lxmartini

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/layer-x/layerx-commons/lxlog"
	"net/http"
	"time"
)

func QuietMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	customMartini := martini.New()
	customMartini.Use(customLogger())
	customMartini.Use(martini.Recovery())
	customMartini.Use(martini.Static("public"))
	customMartini.MapTo(r, (*martini.Routes)(nil))
	customMartini.Action(r.Handle)
	return &martini.ClassicMartini{customMartini, r}
}

func customLogger() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}

		lxlog.Debugf(logrus.Fields{}, fmt.Sprintf("Started %s %s for %s", req.Method, req.URL.Path, addr))

		rw := res.(martini.ResponseWriter)
		c.Next()

		lxlog.Debugf(logrus.Fields{}, fmt.Sprintf("Completed %v %s in %v\n", rw.Status(), http.StatusText(rw.Status()), time.Since(start)))
	}
}