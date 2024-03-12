package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fogleman/gg"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}
func textToImageHandler(w http.ResponseWriter, r *http.Request, fontSize float64, fontPath string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extracting text from the URL path, assuming the URL is like /text-to-image/{text}
	// Split the URL path and extract the text part
	parts := strings.SplitN(r.URL.Path, "/text-to-image/", 2)
	if len(parts) != 2 || parts[1] == "" {
		http.Error(w, "Text must be provided in the URL path", http.StatusBadRequest)
		return
	}
	text := parts[1]

	// Create an image with the text with a transparent background

	const padding = 20
	dc := gg.NewContext(100, 100)
	width, height := dc.MeasureString(text)
	width += 2 * padding  // Add padding to the width
	height += 2 * padding // Add padding to the height

	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dc = gg.NewContext(int(width), int(height)) // Adding a little padding
	dc.SetRGBA(0, 0, 0, 0)                      // Transparent background
	dc.Clear()
	dc.SetRGB(0, 0, 0) // Black text
	dc.DrawStringAnchored(text, width/2, height/2, 0.5, 0.5)
	dc.Clip()

	// Write the image to the response
	w.Header().Set("Content-Type", "image/png")
	dc.EncodePNG(w)
}

func main() {
	// Simple dynamic route handling
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8080" // Default port
	}

	fontSizeStr, exists := os.LookupEnv("FONTSIZE")
	if !exists {
		fontSizeStr = "24" // Default font size
	}
	fontSize, _ := strconv.ParseFloat(fontSizeStr, 64)

	fontPath, exists := os.LookupEnv("FONTPATH")
	if !exists {
		fontPath = "./HackNerdFont-Bold.ttf" // Default font path
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/text-to-image/") {
			textToImageHandler(w, r, fontSize, fontPath)
		} else {
			http.NotFound(w, r)
		}
	})

	http.ListenAndServe(":"+port, nil)
}
