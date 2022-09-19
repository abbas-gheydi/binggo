package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/reujab/wallpaper"
)

const (
	bingAddress    string = "https://www.bing.com/"
	bingApiAddress string = "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=en-US"
)

func main() {

	imageUrl, imageName, err := get_image_name_url()
	if err != nil {
		log.Fatal(err)

	}

	log.Println(imageName)

	if strings.Contains(os.Getenv("XDG_DATA_DIRS"), "xfce") {
		wallpaper.Desktop = "XFCE"

	}
	err = wallpaper.SetFromURL(imageUrl)
	if err != nil {
		log.Fatal(err)

	}

}

func downloadJson() (json []byte, err error) {
	curl := http.Client{Timeout: time.Second * 10}
	resp, err := curl.Get(bingApiAddress)
	if err != nil {
		log.Println("could not get url")
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
		log.Println(err)
		return
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Println(err)
		return
	}

	return

}
