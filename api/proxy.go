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
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/coast-team/mute-auth-proxy/helper"
)

// BotStorageReverseProxy is a structure that contains the needed information for the proxy
type BotStorageReverseProxy struct {
	target         *url.URL               // URL of the target to which the requests are proxied
	LocationPrefix string                 // The listening location path
	proxy          *httputil.ReverseProxy // Actual http reverse proxy
}

// New creates a BotStorageReverseProxy given the ListeningPath and the target
func New(listeningPath string, target string) *BotStorageReverseProxy {
	url, _ := url.Parse(target)
	return &BotStorageReverseProxy{target: url, LocationPrefix: listeningPath, proxy: httputil.NewSingleHostReverseProxy(url)}
}

// Handle checks the JWT and proxies the request to the botstorage
func (p *BotStorageReverseProxy) Handle(w http.ResponseWriter, r *http.Request) error {
	token, err := helper.ExtractJWT(r)
	if err != nil {
		err = helper.IsJWTValid(token, err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return fmt.Errorf("Couldn't extract or validate the JWT.\nError was: %s", err)
	}
	URI := strings.TrimPrefix(r.RequestURI, p.LocationPrefix)
	log.Printf("Botstorage Proxy : URI -> %s", URI)
	p.updateRequestURL(URI)
	p.proxy.ServeHTTP(w, r)
	return nil
}

// updateRequestURL updates the URL in the request by removing the location prefix
// e.g. the auth proxy listens to a /botstorage route, and if a /botstorage/name requests arrives it, this function will rewrite this request's URLPath to /name
func (p *BotStorageReverseProxy) updateRequestURL(newURL string) {
	p.proxy.Director = func(outReq *http.Request) {
		outReq.URL.Scheme = "http"
		outReq.URL.Host = p.target.Host
		outReq.Host = p.target.Host
		outReq.URL.Path = newURL
	}
}
