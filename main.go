package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
)

type Handler struct {
	config *Config
}

var (
	pathRegexp = regexp.MustCompile(`^/([^\/]+)/([^\/]+)/?$`)
	commands   = map[string][]string{
		"up":   {"up", "-d"},
		"stop": {"stop"},
		"pull": {"pull"},
		"rm":   {"rm", "-f"},
	}
	sequences = map[string][]string{
		"start":   {"pull", "up"},
		"stop":    {"stop", "rm"},
		"restart": {"stop", "rm", "pull", "up"},
	}
)

func main() {
	configPath := flag.String("config", "./config.yml", "path to config file")
	flag.Parse()

	config, err := NewConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Loaded config: %#v\n", config)

	h := Handler{config}
	http.HandleFunc("/", h.Handle)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil))
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	match := pathRegexp.FindStringSubmatch(r.URL.Path)
	if match == nil {
		http.NotFound(w, r)
		return
	}

	service := match[1]
	if _, ok := h.config.Services[service]; !ok {
		http.NotFound(w, r)
		return
	}

	sequence := match[2]
	if _, ok := sequences[sequence]; !ok {
		http.NotFound(w, r)
		return
	}

	var err error
	for _, command := range sequences[sequence] {
		cmd := exec.Command(h.config.Binary, commands[command]...)
		cmd.Dir = h.config.Services[service].Path

		if err = cmd.Run(); err != nil {
			break
		}
	}

	if err != nil {
		log.Println("error:", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.Write([]byte("ok"))
	}
}
