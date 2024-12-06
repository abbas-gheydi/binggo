package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	bingAddress    string = "https://www.bing.com/"
	bingApiAddress string = bingAddress + "HPImageArchive.aspx?format=js&idx=0&n=1&mkt=en-US"
	Green                 = "\033[32m"
	Red                   = "\033[31m"
	ResetColor            = "\033[0m"
)

func main() {

	imageUrl, imageName, err := get_image_name_url()
	if err != nil {
		logRedAndExit(err.Error())

	}

	logGreen("this image will be sets as s wallpaper: " + imageName)

	if strings.Contains(os.Getenv("XDG_DATA_DIRS"), "xfce") {
		wallpaper.Desktop = "XFCE"

	}
	// Check if the OS is Darwin (macOS)
	if runtime.GOOS == "darwin" {
		if err := exec.Command("osascript", "-e", `tell application "System Events" to tell every desktop to set picture to "DEFAULT"`).Run(); err != nil {
			logRedAndExit("could not set wallpaper to default" + err.Error())
		}

	}

	err = wallpaper.SetFromURL(imageUrl)
	if err != nil {
		logRedAndExit(err.Error())
	}
	logGreen("the wallpaper is set to: " + imageName)

}

func logRedAndExit(message string) {
	log.Println(Red + message + ResetColor)
	os.Exit(1)
}

func logGreen(message string) {
	log.Println(Green + message + ResetColor)
}

func downloadJson() (json []byte, err error) {
	curl := http.Client{Timeout: time.Second * 10}
	resp, err := curl.Get(bingApiAddress)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	jsonfile, err := ioutil.ReadAll(resp.Body)
	return jsonfile, err

}

func get_image_name_url() (image_url string, image_name string, err error) {
	res, err := Bing_Response{}.receive()
	if err != nil {
		return
	}
	image_url = fmt.Sprint(bingAddress, res.Images[0].URL)
	image_name = fmt.Sprint(res.Images[0].Title, ".jpg")

	return
}

type images struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type Bing_Response struct {
	Images []images `json:"images"`
}

func (b Bing_Response) receive() (res Bing_Response, err error) {

	data, err := downloadJson()
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return
	}

	return

}
