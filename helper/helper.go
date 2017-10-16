package helper

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

var Secret = []byte("random")

var allowedOrigins = []string{"127.0.0.1", "localhost"}

func IsOriginAllowed(origin string) bool {
	allowedOriginsJoined := strings.Join(allowedOrigins, "|")
	var pattern = regexp.MustCompile(fmt.Sprintf(`(https?:\/\/)(%s)(:)([0-9]+)`, allowedOriginsJoined))

	return pattern.MatchString(origin)
}

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

func ExtractJWT(r *http.Request) (*jwt.Token, error) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, jwtKeyFunc)
	return token, err
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return Secret, nil
}

func IsJWTValid(token *jwt.Token, tokenError error) error {
	var msg string
	if token.Valid {
		log.Println("JWT is Valid")
		return nil
	} else if ve, ok := tokenError.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			msg = "That's not even a jwt"
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			msg = "Token is either expired or not active yet"
		} else {
			msg = "Couldn't handle this token"
		}
	} else {
		msg = "Couldn't handle this token"
	}
	return fmt.Errorf(fmt.Sprintf("Error with JWT - %s\nError: %s", msg, tokenError))
}
