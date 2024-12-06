package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/reujab/wallpaper"
)

const (
	bingAPIURL = "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=en-US"
	greenColor = "\033[32m"
	redColor   = "\033[31m"
	resetColor = "\033[0m"
)

func main() {
	// Fetch the image URL and name
	imageURL, imageName, err := getImageNameURL()
	if err != nil {
		logAndExitWithRed(err.Error())
	}

	logWithGreen(fmt.Sprintf("This image will be set as wallpaper: %s", imageName))

	// Check if the environment is XFCE and adjust wallpaper setting accordingly
	if strings.Contains(os.Getenv("XDG_DATA_DIRS"), "xfce") {
		wallpaper.Desktop = "XFCE"
	}

	// If the OS is macOS (Darwin), reset the wallpaper first
	if runtime.GOOS == "darwin" {
		if err := resetMacWallpaper(); err != nil {
			logAndExitWithRed(fmt.Sprintf("Could not reset macOS wallpaper: %v", err))
		}
	}

	// Set the wallpaper from the image URL
	if err := wallpaper.SetFromURL(imageURL); err != nil {
		logAndExitWithRed(fmt.Sprintf("Error setting wallpaper: %v", err))
	}

	logWithGreen(fmt.Sprintf("Wallpaper set to: %s", imageName))
}

// logAndExitWithRed logs a message in red and exits the program.
func logAndExitWithRed(message string) {
	log.Println(redColor + message + resetColor)
	os.Exit(1)
}

// logWithGreen logs a message in green.
func logWithGreen(message string) {
	log.Println(greenColor + message + resetColor)
}

// fetchJSON makes an HTTP GET request to the Bing API and returns the response JSON.
func fetchJSON() ([]byte, error) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(bingAPIURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Bing API: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

// getImageNameURL fetches the image URL and name from Bing's API response.
func getImageNameURL() (string, string, error) {
	response, err := getBingResponse()
	if err != nil {
		return "", "", err
	}

	imageURL := fmt.Sprintf("https://www.bing.com%s", response.Images[0].URL)
	imageName := fmt.Sprintf("%s.jpg", response.Images[0].Title)

	return imageURL, imageName, nil
}

// getBingResponse fetches and unmarshals the response from Bing's API.
func getBingResponse() (BingResponse, error) {
	data, err := fetchJSON()
	if err != nil {
		return BingResponse{}, err
	}

	var response BingResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return BingResponse{}, fmt.Errorf("failed to unmarshal Bing API response: %w", err)
	}

	return response, nil
}

// resetMacWallpaper resets the macOS wallpaper to the default.
func resetMacWallpaper() error {
	cmd := exec.Command("osascript", "-e", `tell application "System Events" to tell every desktop to set picture to "DEFAULT"`)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute AppleScript command: %w", err)
	}
	return nil
}

// BingResponse represents the structure of the JSON response from Bing's API.
type BingResponse struct {
	Images []Image `json:"images"`
}

// Image represents each image object in the Bing API response.
type Image struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}
