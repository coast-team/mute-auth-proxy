package helper

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var allowedOrigins = []string{"127.0.0.1", "localhost"}

// IsOriginAllowed returns True if the current Origin is one of the allowed one
func IsOriginAllowed(origin string) bool {
	allowedOriginsJoined := strings.Join(allowedOrigins, "|")
	var pattern = regexp.MustCompile(fmt.Sprintf(`(https?:\/\/)(%s)(:)([0-9]+)`, allowedOriginsJoined))

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
