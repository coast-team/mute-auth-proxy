// Copyright 2017 Jean-Philippe Eisenbarth
//
// This file is part of Mute Authentication Proxy.
//
// Foobar is free software: you can redistribute it and/or modify
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
// along with Foobar. See the file COPYING.  If not, see <http://www.gnu.org/licenses/>.

package helper

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var allowedOrigins = []string{"127.0.0.1", "localhost", "www.coedit.re"}

// IsOriginAllowed returns True if the current Origin is one of the allowed one
func IsOriginAllowed(origin string) bool {
	log.Printf(origin)
	allowedOriginsJoined := strings.Join(allowedOrigins, "|")
	var pattern = regexp.MustCompile(fmt.Sprintf(`(https?:\/\/)(%s)(:[0-9]+)?`, allowedOriginsJoined))

	return pattern.MatchString(origin)
}

// SetHeader adds to the header the 'good' value for CORS
func SetHeader(w http.ResponseWriter, r *http.Request) {
	if IsOriginAllowed(r.Header.Get("Origin")) {
		log.Printf("Origin %s allowed\n", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin")) // IMPORTANT
		w.Header().Set("Vary", "Origin, Access-Control-Request-Headers")
		w.Header().Set("Access-Control-Allow-Methods", "GET,HEAD,PUT,PATCH,POST,DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "content-type, authorization") // IMPORTANT !
		w.Header().Set("Connection", "keep-alive")
	}
}
