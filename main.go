package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
)

var (
	// Configuration variables
	listenPort = getEnv("PORT", "8080")
	fontSize   = parseFloatEnv("FONT_SIZE", 24)
	fontPath   = getEnv("FONT_PATH", "./HackNerdFont-Bold.ttf")
)

func main() {
	r := gin.Default()

	r.GET("/text-to-image/:text", textToImageHandler)
	// r.GET("/font-as-bytes", fontAsBytesHandler)

	//logMessage(fmt.Sprintf("Listening on port %s", listenPort))
	if err := r.Run(":" + listenPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func textToImageHandler(c *gin.Context) {
	text := c.Param("text")
	//logRequest(c.Request, text)
	safeText := url.PathEscape(strings.ReplaceAll(text, "/", "_")) // Basic sanitization for filename
	cleantext, err := url.QueryUnescape(safeText)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("Input text is : %v", safeText)
	log.Printf("Cleaned text is : %v", cleantext)
	dc := gg.NewContext(100, 100)
	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		c.String(http.StatusInternalServerError, "Failed to load font face")
		return
	}

	width, height := dc.MeasureString(cleantext)
	dc = gg.NewContext(int(width)+10, int(height)+10) // Add some padding
	dc.SetRGBA(0, 0, 0, 0)                            // Transparent background
	dc.Clear()
	dc.SetRGB(0, 0, 0) // Black text
	dc.DrawStringAnchored(cleantext, width/2, height/2, 0.5, 0.5)
	dc.Clip()

	//logResponse(text, int(width)+10, int(height)+10)

	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf); err != nil {
		c.String(http.StatusInternalServerError, "Failed to encode image")
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.png", cleantext))
	c.Data(http.StatusOK, "image/png", buf.Bytes())
}

// func logRequest(r *http.Request, text string) {
// 	logMessage(fmt.Sprintf("Request: %s, Path: %s", text, r.URL.Path))
// }

// func logResponse(text string, width, height int) {
// 	logMessage(fmt.Sprintf("Response: %s, Image Size: %dx%d", text, width, height))
// }

// func logMessage(message string) {
// 	log.Printf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), message)
// }

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func parseFloatEnv(key string, fallback float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return fallback
}
