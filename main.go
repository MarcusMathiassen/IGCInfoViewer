package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/marni/goigc"

	"github.com/gin-gonic/gin"
)

type TrackInfo struct {
	pilot       string
	glider      string
	gliderID    string
	trackLength string
}

var trackInfos []TrackInfo
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
			id := len(tracks)
			tracks = append(tracks, url)
			c.JSON(http.StatusOK, gin.H{"id": id})
		})

		api.GET("/igc/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				panic(err)
			}
			if id >= len(tracks) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			trackURL := tracks[id]
			track, err := igc.ParseLocation(trackURL)

			trackInfo := TrackInfo{
				pilot:       track.Pilot,
				glider:      track.GliderType,
				gliderID:    track.GliderID,
				trackLength: track.UniqueID,
			}
			// trackInfos[track.UniqueID] = trackInfo
			trackInfos = append(trackInfos, trackInfo)

			c.JSON(http.StatusOK, gin.H{
				"H_date":       track.Header,
				"pilot":        trackInfo.pilot,
				"glider":       trackInfo.glider,
				"glider_id":    trackInfo.gliderID,
				"track_length": trackInfo.trackLength,
			})
		})
		api.GET("/igc", func(c *gin.Context) {
			// We return an empty array if there are no tracks yet.
			// Gin turns the empty 'tracks' array into 'null'
			//  so we create a temporary '[]int' to satisfy the requirements.
			if len(tracks) == 0 {
				c.JSON(http.StatusOK, []int{})
				return
			}
			c.JSON(http.StatusOK, tracks)
		})
	}
	r.Run(":" + port)
}
