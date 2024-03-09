package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*		if strings.HasPrefix(r.URL.Path, "/docs") {
					http.ServeFile(w, r, "/docs/swagger.json")
					return
				}
				if strings.HasPrefix(r.URL.Path, "/swagger") {
					http.ServeFile(w, r, "/docs/swagger.json")
					return
				}*/
		if strings.HasPrefix(r.URL.Path, "/swagger") {
			next.ServeHTTP(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/mycustompath") {
			next.ServeHTTP(w, r)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/api") {
			target := &url.URL{
				Scheme: "http",
				Host:   rp.host + ":" + rp.port,
			}

			proxy := httputil.NewSingleHostReverseProxy(target)
			originalDirector := proxy.Director
			proxy.Director = func(req *http.Request) {
				originalDirector(req)
				req.URL.Host = target.Host
				req.URL.Scheme = target.Scheme
				req.Host = target.Host
			}

			proxy.ModifyResponse = func(response *http.Response) error {
				if response.StatusCode == http.StatusNotFound {
					http.Redirect(w, r, "http://hugo:1313/", http.StatusFound)
					return nil
				}
				return nil
			}

			proxy.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
