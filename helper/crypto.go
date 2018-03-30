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
	"crypto/rand"
	"log"
)

type secret struct {
	secret     []byte
	setCounter int
}

var sec = &secret{setCounter: 0}

// GetSecret returns the secret that should be used for signing JWT
func GetSecret() []byte {
	if sec.setCounter != 1 {
		log.Fatal("The secret has not yet being set ...")
	}
	return sec.secret
}

// GenerateSecret sets the secret by generating a new one
// The secret can only be set once...
func GenerateSecret() {
	if sec.setCounter != 0 {
		log.Println("The secret has already being set ...")
	} else {
		sec.secret = GenerateRandomBytes()
	}
}

// SetSecret sets the secret given b an array of bytes read from a file for example
func SetSecret(b []byte) {
	if sec.setCounter != 0 {
		log.Println("The secret has already being set ...")
	} else {
		sec.secret = b
		sec.setCounter++
	}
}

// GenerateRandomBytes generated a []byte containing n secure random numbers
func GenerateRandomBytes() []byte {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		log.Fatal(err)
	}
	return b
}
