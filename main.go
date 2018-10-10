package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Uptime struct {
	uptime  string
	info    string
	version string
}

func handleBadRequest(c *gin.Context) {
	fmt.Println("handleBadRequest ...")
}

func handleRequest(c *gin.Context) {
	fmt.Println("handleRequest ...")
}

func main() {

	// s := "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	// track, err := igc.ParseLocation(s)
	// if err != nil {
	// 	fmt.Errorf("Problem reading the track %s", err)
	// }
	// fmt.Printf("Pilot: %s, gliderType: %s, date: %s\n", track.Pilot, track.GliderType, track.Date.String())

	r := gin.Default()
	r.GET("/", handleBadRequest)
	r.GET("/igcinfo/api", handleRequest)
	r.Run()
}
