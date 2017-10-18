package api

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"github.com/coast-team/mute-auth-proxy/config"
	"github.com/coast-team/mute-auth-proxy/helper"
)

// MakeConiksProxyHandler is the handler for the route that proxies a Coniks request
func MakeConiksProxyHandler(conf *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		helper.SetHeader(w, r)
		switch r.Method {
		case "POST":
			handleConiksProxy(w, r, conf)
		default:
			log.Printf("Method : %s", r.Method)
			log.Printf("Request: %v", r.Body)
			w.WriteHeader(204)
			return
		}
	}
}

func handleConiksProxy(w http.ResponseWriter, r *http.Request, conf *config.Config) {
	token, err := helper.ExtractJWT(r)
	if err != nil {
		err = helper.IsJWTValid(token, err)
		log.Printf("JWT error: %s", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Post(conf.ConiksServerAddr, "application/json", r.Body)
	if err != nil {
		log.Fatalf("Coniksserver request failed: %s", err.Error())
	}
	io.Copy(w, res.Body)
}
