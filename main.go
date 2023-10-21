package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/net/webdav"
	"gopkg.in/yaml.v3"
)

type Account struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dir      string `yaml:"dir"`
	handler  *webdav.Handler
}

type Config struct {
	Port     int       `yaml:"port"`
	Accounts []Account `yaml:"accounts"`
}

func main() {
	configPath := flag.String("c", "config.yml", "specify configuration file")
	flag.Parse()

	// read config file
	file, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	// parse config file
	config := Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatalln(err)
	}

	if len(config.Accounts) == 0 {
		log.Fatalln("accounts is required in config file")
	}

	// init account handler
	cache := make(map[string]Account)
	for i := 0; i < len(config.Accounts); i++ {
		account := config.Accounts[i]

		// check duplicate
		if _, ok := cache[account.User]; ok {
			log.Fatalln(errors.New("duplicate account"))
		}

		// check home dir
		absPath, err := filepath.Abs(account.Dir)
		if err != nil {
			log.Fatalln(err)
		}
		info, err := os.Stat(absPath)
		if (err == nil && !info.IsDir()) || err != nil {
			log.Fatalln(errors.New("path does not exist: " + absPath))
		}

		account.handler = &webdav.Handler{
			Prefix:     "/",
			FileSystem: webdav.Dir(absPath),
			LockSystem: webdav.NewMemLS(),
		}
		cache[account.User] = account
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// basic auth
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		account, exist := cache[username]
		if !exist || username != account.User || password != account.Password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		account.handler.ServeHTTP(w, r)
	})

	go func() {
		time.Sleep(time.Second)
		fmt.Println("server is running on port", config.Port)
	}()

	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), mux)
	if err != nil {
		log.Fatalln(err)
	}
}
