package main

import (
	"log"
	"net/http"
	"os"
	"text/template"
)

type job struct {
	Name           string
	Status         string
	StartTime      string
	CompletionTime string
	Duration       string
	Logs           string
}

type JobList struct {
	Namespace    string
	Subscription string
	Jobs         []job
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("jobs.tmpl"))
	tmpl.Execute(w, jobs())
}

func main() {
	log.Print("Starting monitor")

	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Service Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func jobs() JobList {

	jobs := []job{{
		Name:           "engine-123",
		Status:         "Running",
		StartTime:      "123",
		CompletionTime: "456",
		Duration:       "99 s",
		Logs:           "link"}}
	return JobList{
		Namespace:    "namesapce",
		Subscription: "sub",
		Jobs:         jobs}
}
