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

package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config represents the structure containing the information from the config file
type Config struct {
	Port             int
	ConiksServerAddr string      `toml:"coniksserver_addr"`
	KeyServerPath    string      `toml:"keyserver_path"`
	BotStorageAddr   string      `toml:"botstorage_addr"`
	AllowedOrigins   []string    `toml:"allowed_origins"`
	OauthPrefs       OauthConfig `toml:"oauth"`
}

func (conf Config) String() string {
	return fmt.Sprintf("Config:\n  Port: %d\n  Coniks server addr: %s\n  KeyServer path: %s\n  BotStorage addr: %s\n  Allowed origins: %s\n  %s", conf.Port, conf.ConiksServerAddr, conf.KeyServerPath, conf.BotStorageAddr, conf.AllowedOrigins, conf.OauthPrefs)
}

type OauthConfig struct {
	GooglePrefs ProviderPrefs `toml:"google"`
	GithubPrefs ProviderPrefs `toml:"github"`
}

func (conf OauthConfig) String() string {
	return fmt.Sprintf("Oauth Config:\n    Google Preferences:\n      %s\n    Github Preferences:\n      %s", conf.GooglePrefs, conf.GithubPrefs)
}

type ProviderPrefs struct {
	ClientSecret string `toml:"client_secret"`
}

func (conf ProviderPrefs) String() string {
	return fmt.Sprintf("Client Secret: %s", conf.ClientSecret)
}

// LoadConfig loads and parses the information from the config file and fill the Config struct
func LoadConfig(file string) (*Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return nil, fmt.Errorf("Failed to load config: %v", err)
	}

	return &conf, nil
}
