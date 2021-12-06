package main

import (
	"net/url"
	"testing"
)

func Test_get_image_name_url(t *testing.T) {
	image, _, err := get_image_name_url()
	if err != nil {
		t.Error(err)
	}
	//test if url is valied
	if _, Err := url.ParseRequestURI(image); Err != nil {
		t.Error(Err)

	}
}
func Test_Bing_Response(t *testing.T) {
	bing := Bing_Response{}
	rec, err := bing.receive()
	if rec.Images == nil {
		t.Error("image is empty")
	}
	if err != nil {
		t.Error(err)
	}
}
func Test_getRawJson(t *testing.T) {
	data, err := downloadJson()
	if data == nil {
		t.Error("empty data")
	}
	if err != nil {
		t.Error(err)
	}

}
