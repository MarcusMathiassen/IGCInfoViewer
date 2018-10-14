package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marni/goigc"
)

type TrackInfo struct {
	ID          int
	TrackLength float64
	Pilot       string `form:"pilot" json:"pilot" binding: required`
	Glider      string `form:"glider" json:"glider" binding: required`
	GliderID    string `form:"glider_id" json:"glider_id" binding: required`
	HDate       string `form:"H_date" json:"H_date" binding: required`
}

var trackInfos []TrackInfo
var startTime = time.Now()

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
func distanceInMetersBetweenGPSCoords(lat1, lon1, lat2, lon2 float64) float64 {
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

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

func (t TrackInfo) getField(fieldName string) string {
	switch fieldName {
	case "pilot":
		return t.Pilot
	case "glider":
		return t.Glider
	case "glider_id":
		return t.GliderID
	case "H_date":
		return t.HDate
	case "track_length":
		return strconv.FormatFloat(t.TrackLength, 'f', 6, 64)
	}
	return "Unknown"
}

func main() {

	router := gin.Default()

	port := os.Getenv("PORT")
	if port == "" { // for running locally
		port = "8080"
	}

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

			track, err := igc.ParseLocation(url)
			if err != nil {
				log.Fatal(err)
			}

			points := track.Points
			trackLength := 0.0
			for i := 1; i < len(points); i++ {
				trackLength += points[i-1].Distance(points[i])
			}

			id := len(trackInfos)
			trackInfo := TrackInfo{
				ID:          id,
				TrackLength: trackLength,
				Pilot:       track.Pilot,
				Glider:      track.GliderType,
				GliderID:    track.GliderID,
				HDate:       track.Header.Date.String(),
			}

			trackInfos = append(trackInfos, trackInfo)
			c.JSON(http.StatusOK, gin.H{"id": id})
		})

		api.GET("/igc", func(c *gin.Context) {
			ids := make([]int, len(trackInfos))
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
			if id >= len(trackInfos) {
				c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
				return
			}

			trackInfo := trackInfos[id]

			c.JSON(http.StatusOK, gin.H{
				"H_date":       trackInfo.HDate,
				"pilot":        trackInfo.Pilot,
				"glider":       trackInfo.Glider,
				"glider_id":    trackInfo.GliderID,
				"track_length": trackInfo.TrackLength,
			})
		})

		api.GET("/igc/:id/:field", func(c *gin.Context) {
			field := c.Param("field")
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if id >= len(trackInfos) {
				c.Status(http.StatusNotFound)
				return
			}
			trackInfo := trackInfos[id]
			c.String(http.StatusOK, trackInfo.getField(field))
		})
	}

	router.Run(":" + port)
}
