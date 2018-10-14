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

func getTrackByID(id int) {

}

func main() {

	port := os.Getenv("PORT")
	if port == "" { // for running locally
		port = "8080"
	}

	router := gin.Default()

	router.GET("/igcinfo/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"uptime":  getUptime(),
			"info":    "Service for IGC tracks.",
			"version": "v1",
		})
	})

	api := router.Group("/igcinfo/api")
	{
		api.POST("/igc", func(c *gin.Context) {
			url, exists := c.GetPostForm("url")
			if !exists {
				c.JSON(http.StatusBadRequest, gin.H{"error": "missing key 'url'"})
				return
			}
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
				c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
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

		api.GET("/igc/:id/:field", func(c *gin.Context) {
			field := c.Param("field")
			id, err := strconv.Atoi(c.Param("id"))
			fmt.Printf("id: %d  field: %s", id, field)
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
			switch field {
			case "H_date":
				c.JSON(http.StatusOK, gin.H{"H_date": track.Header})
			case "pilot":
				c.JSON(http.StatusOK, gin.H{"pilot": track.Pilot})
			case "glider":
				c.JSON(http.StatusOK, gin.H{"glider": track.GliderType})
			case "glider_id":
				c.JSON(http.StatusOK, gin.H{"glider_id": track.GliderID})
			case "track_length":
				c.JSON(http.StatusOK, gin.H{"track_length": track.UniqueID})
			}
		})
	}

	router.Run(":" + port)
}
