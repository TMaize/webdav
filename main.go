package main

import (
	"errors"
	"flag"
	"fmt"
	"golang.org/x/net/webdav"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Account struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dir      string `yaml:"dir"`
}

type Config struct {
	Port     int       `yaml:"port"`
	Accounts []Account `yaml:"accounts"`
}

func checkName(name string) error {
	matched, err := regexp.MatchString("^[0-9a-zA-Z]+$", name)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("user need match [0-9a-zA-Z]")
	}
	return nil
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
	handlers := make(map[string]webdav.Handler, 0)
	for _, account := range config.Accounts {
		err := checkName(account.User)
		if err != nil {
			log.Fatalln(err)
		}
		handlers[account.User] = webdav.Handler{
			Prefix:     "/user/" + account.User,
			FileSystem: webdav.Dir(account.Dir),
			LockSystem: webdav.NewMemLS(),
		}
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	})

	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		arr := strings.Split(r.URL.Path, "/")
		user := arr[2]

		var account *Account
		for _, item := range config.Accounts {
			if item.User == user {
				account = &item
			}
		}

		if account == nil {
			w.WriteHeader(404)
			return
		}

		// basic auth
		username, password, ok := r.BasicAuth()
		if account.Password != "" && (!ok || username != account.User || password != account.Password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler := handlers[account.User]
		handler.ServeHTTP(w, r)
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
