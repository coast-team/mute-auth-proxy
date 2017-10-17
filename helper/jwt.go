package helper

import (
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

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
