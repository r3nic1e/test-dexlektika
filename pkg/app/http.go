package app

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/r3nic1e/test-dexlektika/pkg/models"
)

func (app *App) StartServer(addr string) error {
	app.mux.HandleFunc("/", app.squareHandler)
	app.mux.HandleFunc("/healthy", app.healthyHandler)
	app.mux.HandleFunc("/ready", app.readyHandler)
	// I don't support using this word but it is required to call the method like that
	app.mux.HandleFunc("/blacklisted", app.denylistedHandler)
	app.mux.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addr, app.mux)
}

func (app *App) denylistedHandler(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		log.Printf("Failed to parse remote addr %s: %v", request.RemoteAddr, err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	result := app.db.Create(models.DenylistedIP{
		IP:       ip,
		HTTPPath: request.RequestURI,
	})
	if result.Error != nil {
		log.Printf("Failed to create denylisted IP: %v", result.Error)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(444)
}

func (app *App) isDenylisted(request *http.Request) bool {
	ip, _, err := net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		log.Printf("Failed to parse remote addr %s: %v. Denying request", request.RemoteAddr, err)
		return true
	}

	var count int64
	app.db.Select(&models.DenylistedIP{IP: ip}).Count(&count)

	return count > 0
}

func (app *App) squareHandler(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	values := request.URL.Query()
	if !values.Has("n") {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	nStr := values.Get("n")
	n, err := strconv.Atoi(nStr)
	if err != nil {
		log.Printf("Failed to parse %q as int: %v", nStr, err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	if app.isDenylisted(request) {
		response.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(response, "IP is in denylist")
		return
	}

	_, err = fmt.Fprint(response, n*n)
	if err != nil {
		log.Printf("Failed to write HTTP response: %v", err)
		response.WriteHeader(500)
	}
}

func (app *App) healthyHandler(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	response.WriteHeader(http.StatusOK)
}

func (app *App) readyHandler(response http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()

	response.WriteHeader(http.StatusOK)
}
