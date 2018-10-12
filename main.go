package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

func fmtTimeDifference(t time.Time) string {
	diff := time.Since(t)
	return diff.String()
}
func getUptime() string {
	return fmtTimeDifference(startTime)
}

func handleBadRequest(c *gin.Context) {
	fmt.Println("handleBadRequest ...")
}

func handleRequest(c *gin.Context) {
	url, err := readAll(c.Request.Body)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"url": url})
}
func main() {
	port := os.Getenv("PORT")
	if port == "" { // for running locally
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
		// c.Writer.Header().Set("Content-Type", "application/json")
		c.IndentedJSON(http.StatusOK, gin.H{
			"uptime":  getUptime(),
			"info":    "Service for IGC tracks.",
			"version": "v1",
		})
	})
	r.POST("/igcinfo/api/igc", handleRequest)
	r.Run(":" + port)
}
