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

package helper

import (
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// ExtractJWT parses the http.Request and extract the JWT from it.
// It also check that the token signature is correct.
func ExtractJWT(r *http.Request) (*jwt.Token, error) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, jwtKeyFunc)
	return token, err
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	return GetSecret(), nil
}

// IsJWTValid checks that the token is a well formed and not expired JWT
func IsJWTValid(token *jwt.Token, tokenError error) error {
	var msg string
	if token != nil && token.Valid {
		log.Println("JWT is Valid")
		return nil
	} else if ve, ok := tokenError.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			msg = "That's not even a JWT"
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			msg = "Token is either expired or not active yet"
		} else {
			msg = "Couldn't handle this token"
		}
	} else {
		msg = "Couldn't handle this token"
	}
	return fmt.Errorf("%s: %s", msg, tokenError)
}

func GenerateJWT() *jwt.Token {
	return jwt.New(jwt.SigningMethodHS256)
}

func GetSignedString(token *jwt.Token) (string, error) {
	return token.SignedString(GetSecret())
}
