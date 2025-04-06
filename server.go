package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var Version string 

type App struct {
	Port string
}

func (a *App) Start() {
	addr := fmt.Sprintf(":%s", a.Port)
	log.Printf("Starting app on %s", addr)
	http.HandleFunc("/", index)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func env(key, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return val
}

func main() {
	server := App{
		Port: env("PORT", "8080"),
	}
	server.Start()
}

func index(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Hostname string
		IP       string
		Version  string
	}{
		Hostname: getHostname(),
		IP:       getIP(),
		Version:  getVersion(),
	}
	renderTemplate(w, "index.html", data)
}

func renderTemplate(w http.ResponseWriter, templateName string, data interface{}) {
    tmpl, err := template.ParseFiles("/usr/share/nginx/html/index.html")
    if err != nil {
        log.Printf("Template parsing error: %v", err)
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
        return
    }
    if err := tmpl.Execute(w, data); err != nil {
        log.Printf("Template execution error: %v", err)
        http.Error(w, "Error rendering template", http.StatusInternalServerError)
    }
}


func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func getIP() string {
	return "127.0.0.1"
}



func getVersion() string {
	if Version != "" {
		return Version
	}
	return "1.0.0"
}