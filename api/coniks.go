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
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/coast-team/mute-auth-proxy/config"
	"github.com/coast-team/mute-auth-proxy/helper"
)

// MakeConiksProxyHandler is the handler for the route that proxies a Coniks request
func MakeConiksProxyHandler(conf *config.Config) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handleConiksProxy(w, r, conf)
		if err != nil {
			log.Printf("Coniks proxy err: %s\n", err)
		}
	}
}

func handleConiksProxy(w http.ResponseWriter, r *http.Request, conf *config.Config) error {
	token, err := helper.ExtractJWT(r)
	if err != nil {
		err = helper.IsJWTValid(token, err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return fmt.Errorf("Couldn't extract or validate the JWT.\nError was: %s", err)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("Couldn't read request's body.\nError was: %s", err)
	}

	tlsConf := &tls.Config{InsecureSkipVerify: true}
	u, _ := url.Parse(conf.ConiksServerAddr)
	conn, err := tls.Dial(u.Scheme, u.Host, tlsConf)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("Couldn't establish connection to ConiksServer.\nError was: %s", err)
	}
	defer conn.Close()

	_, err = conn.Write(body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("Communication to ConiksServer failed. Tried to send:\n%s\nError was: %s", body, err)
	}
	conn.CloseWrite() // writes EOF

	var buf bytes.Buffer
	_, err = io.Copy(&buf, conn)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return fmt.Errorf("Couldn't send ConiksServer's response to ConiksClient.\nError was: %s", err)
	}

	w.Write(buf.Bytes())
	return nil
}
