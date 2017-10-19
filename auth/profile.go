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

package auth

import "fmt"

type Token struct {
	AccessToken string `json:"access_token"`
}

// Profile is the struct representing a User
// For now if a value equals to "" it means that the identity provider did not provide this value.
type Profile struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Picture   string `json:"picture"`
}

// Details returns a string containing all the usefull information about an user.
// Used for logging
func (p Profile) Details() string {
	return fmt.Sprintf("Profile:\n\tlogin: %s\n\temail: %s\n\tfullname: %s\n\tavatar: %s", p.UserLogin(), p.Email, p.Name, p.Avatar())
}

// UserLogin returns the login of the user if provided. Otherwise it return the email of the user.
func (p Profile) UserLogin() string {
	if p.Login != "" {
		return p.Login
	}
	return p.Email
}

// Avatar returns the URL (provider host) that points to the avatar of the user
func (p Profile) Avatar() string {
	if p.AvatarURL != "" {
		return p.AvatarURL
	}
	return p.Picture
}
