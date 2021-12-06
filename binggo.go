package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	bingAddress    string = "https://www.bing.com/"
	bingApiAddress string = "https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=en-US"
)

var (
	imageFolder string = func() string {
		home, _ := os.UserHomeDir()
		return fmt.Sprint(home, "/Pictures/binggo")

	}()
)

func init() {
	err := mkdir()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {

	exitOnError := func(err error) {
		log.Println(err)
		os.Exit(1)
	}

	// connect to bing api and get image url
	imageUrl, imageName, err := get_image_name_url()
	if err != nil {
		exitOnError(err)

	}

	err = download(imageName, imageUrl)
	if err != nil {
		exitOnError(err)

	}

	err = setWallpaper(imageName)
	if err != nil {
		exitOnError(err)

	}

}

func setWallpaper(wallpaper string) error {
	app := "pcmanfm"
	arg0 := "--set-wallpaper"
	arg1 := imageFolder + "/" + wallpaper
	cmd := exec.Command(app, arg0, arg1)

	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	if string(stdout) != "" {
		log.Println(string(stdout))
	}
	log.Println("wallpaper changed to", arg1)
	return nil

}

func download(name string, url string) error {

	//download
	httpClient := http.Client{}
	response, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//create empty image file
	imageOndisk := fmt.Sprint(imageFolder, "/", name)
	image, err := os.Create(imageOndisk)
	if err != nil {
		return err
	}

	defer image.Close()

	//save data to disk
	_, err = io.Copy(image, response.Body)
	log.Println("wallpaper is saved to", imageOndisk)

	return err
}

func mkdir() error {

	info, err := os.Stat(imageFolder)
	// try to create folder
	if os.IsNotExist(err) {

		mkdir_err := os.Mkdir(imageFolder, 0755)
		if mkdir_err == nil {
			log.Println(imageFolder, "created")
			return nil
		} else {
			return mkdir_err

		}
	}

	// if there is a file withe same folder name in the path
	if err == nil && !info.IsDir() {
		err = errors.New(fmt.Sprint(imageFolder, " is not a folder"))
		return err
	}

	return err

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

type Bing_Response struct {
	Images []images `json:"images"`
}
type images struct {
	URL   string `json:"url"`
	Title string `json:"title"`
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

func downloadJson() (json []byte, err error) {
	curl := http.Client{Timeout: time.Second * 3}
	resp, err := curl.Get(bingApiAddress)
	if err != nil {
		log.Println("could not get url")
		return
	}
	defer resp.Body.Close()
	jsonfile, err := ioutil.ReadAll(resp.Body)
	return jsonfile, err

}
