package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/marni/goigc"
)

type Uptime struct {
	uptime  string
	info    string
	version string
}

func handleBadRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("handleBadRequest ...")
	fmt.Fprintf(w, "Bad request.\n")
}

func handleRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("handleRequest ...")
	if r.Method == "GET" {
		name := ps.ByName("name")
		switch name {
		case "api":
			u := Uptime{
				uptime:  time.Now().String(),
				info:    "Service for IGC tracks.",
				version: "v1",
			}
			b, err := json.Marshal(u)
			if err != nil {
				log.Fatal(err)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Printf("%s", b)
		}
	}
}
func handleUptime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("handleUptime ...")
	if r.Method == "GET" {
		name := ps.ByName("name")
		switch name {
		case "api":
			u := Uptime{
				uptime:  time.Now().String(),
				info:    "Service for IGC tracks.",
				version: "v1",
			}
			b, err := json.Marshal(u)
			if err != nil {
				log.Fatal(err)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Printf("%s", b)
		}
	}
}
func main() {

	s := "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	track, err := igc.ParseLocation(s)
	if err != nil {
		fmt.Errorf("Problem reading the track %s", err)
	}
	fmt.Printf("Pilot: %s, gliderType: %s, date: %s\n", track.Pilot, track.GliderType, track.Date.String())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := httprouter.New()
	router.GET("/", handleBadRequest)
	router.GET("/igcinfo/api", handleRequest)

	fmt.Println("Listening on port " + port + "...")
	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}
