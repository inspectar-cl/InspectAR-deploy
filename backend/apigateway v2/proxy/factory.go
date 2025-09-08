package proxy

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"apigateway/pkg"
)

// Crea un ReverseProxy con:
// - targetBase: ej "http://localhost:4000"
// - stripPrefix: ej "/api/sensores" (se quita)
// - prependPath: ej "/activo" (se agrega después del strip)
func NewReverseProxy(targetBase, stripPrefix, prependPath string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(targetBase)
	if err != nil {
		log.Fatalf("URL destino inválida (%s): %v", targetBase, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// timeouts sanos
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	original := proxy.Director
	proxy.Director = func(req *http.Request) {
		original(req)

		// set destino
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host

		// reescritura de path
		p := req.URL.Path
		if stripPrefix != "" && strings.HasPrefix(p, stripPrefix) {
			p = strings.TrimPrefix(p, stripPrefix)
			if p == "" {
				p = "/"
			}
		}
		if prependPath != "" {
			p = pkg.JoinPath(prependPath, p)
		}
		req.URL.Path = pkg.JoinPath(targetURL.Path, p)
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[Gateway ERROR] %s %s → %v", r.Method, r.URL.Path, err)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusBadGateway)
		_, _ = rw.Write([]byte(`{"error":true,"message":"Servicio destino no disponible"}`))
	}

	return proxy
}
