package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

func fmtDuration(duration time.Duration) string {
	days := int64(duration.Hours() / 24)
	years := days / 365
	months := years / 12
	hours := int64(math.Mod(duration.Hours(), 24))
	minutes := int64(math.Mod(duration.Minutes(), 60))
	seconds := int64(math.Mod(duration.Seconds(), 60))

	return fmt.Sprintf("P%dY%dM%dDT%dH%dM%dS", years, months, days, hours, minutes, seconds)
}

func getUptime() string {
	return fmtDuration(time.Since(startTime))
}

func handleBadRequest(c *gin.Context) {
	fmt.Println("handleBadRequest ...")
}

func handleRequest(c *gin.Context) {
	c.String(http.StatusOK, "Hello %s", c.Param("name"))
	fmt.Println("handleRequest ...")
}

func main() {

	port := os.Getenv("PORT")
	if port == "" { // ....if heroku didn't give us a port (DEBUG)
		port = "8080"
	}

	// s := "http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc"
	// track, err := igc.ParseLocation(s)
	// if err != nil {
	// 	fmt.Errorf("Problem reading the track %s", err)
	// }
	// fmt.Printf("Pilot: %s, gliderType: %s, date: %s\n", track.Pilot, track.GliderType, track.Date.String())
	r := gin.Default()
	r.GET("/igcinfo/api", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.IndentedJSON(http.StatusOK, gin.H{
			"uptime":  getUptime(),
			"info":    "Service for IGC tracks.",
			"version": "v1",
		})
	})
	r.GET("/igcinfo/api/:name", handleRequest)
	r.Run(":" + port)
}
