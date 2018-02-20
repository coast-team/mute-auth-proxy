// Copyright 2017-2018 Jean-Philippe Eisenbarth
//
// This file is part of Mute Authentication Proxy.
//
// Mute Authentication Proxy is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Mute Authentication Proxy is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Mute Authentication Proxy. See the file COPYING.  If not, see <http://www.gnu.org/licenses/>.

package api

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"

	"github.com/coast-team/mute-auth-proxy/config"
)

// MakeConiksProxyHandler is the handler for the route that proxies a Coniks request
func MakeConiksProxyHandler(conf *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handleConiksProxy(w, r, conf)
		if err != nil {
			log.Println(err)
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
