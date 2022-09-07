package middlewares

import (
	"github.com/felixge/httpsnoop"
	"net/http"
	"tie.prodigy9.co/config"
)

func AddRequestLogging(handler http.Handler, cfg *config.Config) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, r *http.Request) {
		var metrics *httpsnoop.Metrics
		cfg.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.RequestURI)
		defer func() {
			if metrics == nil {
				return
			}

			cfg.Printf("%s %s %s - HTTP %d %s\n",
				r.RemoteAddr, r.Method, r.RequestURI,
				metrics.Code, metrics.Duration)
		}()

		m := httpsnoop.CaptureMetrics(handler, resp, r)
		metrics = &m
	})

}
