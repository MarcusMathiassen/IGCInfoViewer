package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marni/goigc"
)

var tracks []string
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

func main() {

	port := os.Getenv("PORT")
	if port == "" { // for running locally
		port = "8080"
	}

	r := gin.Default()

	r.GET("/igcinfo/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"uptime":  getUptime(),
			"info":    "Service for IGC tracks.",
			"version": "v1",
		})
	})

	api := r.Group("/igcinfo/api")
	{
		api.POST("/igc", func(c *gin.Context) {
			url := c.PostForm("url")
			if filepath.Ext(url) != ".igc" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "not a .igc file"})
				return
			}
			id := len(tracks)
			tracks = append(tracks, url)
			c.JSON(http.StatusOK, gin.H{"id": id})
		})

		api.GET("/igc", func(c *gin.Context) {
			ids := make([]int, len(tracks))
			for i := range ids {
				ids[i] = i
			}
			c.JSON(http.StatusOK, ids)
		})

		api.GET("/igc/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if id >= len(tracks) {
				c.Status(http.StatusNotFound)
				return
			}
			trackURL := tracks[id]
			track, err := igc.ParseLocation(trackURL)

			c.JSON(http.StatusOK, gin.H{
				"H_date":       track.Header,
				"pilot":        track.Pilot,
				"glider":       track.GliderType,
				"glider_id":    track.GliderID,
				"track_length": track.UniqueID,
			})
		})
	}

	r.Run(":" + port)
}
