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
	"io/ioutil"
	"os"

	"github.com/Songmu/prompter"
)

// WriteFile writes buf to a file whose path is indicated by filepath.
func WriteFile(filepath string, buf []byte, perm os.FileMode) (bool, error) {
	if _, err := os.Stat(filepath); !os.IsNotExist(err) {
		if !prompter.YN(fmt.Sprintf("The file %s already exists, do you want to overwrite it ?", filepath), false) {
			return false, nil
		}
	}
	err := ioutil.WriteFile(filepath, buf, perm)
	return true, err
}
