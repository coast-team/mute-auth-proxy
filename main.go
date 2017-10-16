package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/coast-team/mute-auth-proxy/api"
	"github.com/coast-team/mute-auth-proxy/auth"
	"github.com/coast-team/mute-auth-proxy/config"
)

func main() {
	var confFilename string
	flag.StringVar(&confFilename, "c", "config.toml", "The config file to load - default is 'config.toml'")
	flag.StringVar(&confFilename, "config", "config.toml", "The config file to load - default is 'config.toml' (shorthand)")
	flag.Parse()
	conf, err := config.LoadConfig(confFilename)
	if err != nil {
		log.Fatalf("LoadConfig: %s", err)
	}
	log.Printf("Conf: %+v", conf)
	http.HandleFunc("/auth/google", auth.MakeGoogleLoginHandler(conf))
	http.HandleFunc("/auth/github", auth.MakeGithubLoginHandler(conf))
	http.HandleFunc("/coniks", api.MakeConiksProxyHandler(conf))
	err = http.ListenAndServe(fmt.Sprintf(":%d", conf.Port), nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}
