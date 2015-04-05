package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"text/template"
	"time"
)

var (
	//For some reason soundcloud dosen't append this id on stream_url..
	clientId = flag.String("cid", "", "Client id")

	//Secret key is added on stream_url \o/
	//secret_key = flag.String("apikey", "", "Secret key")

	url    = flag.String("url", "", "Url to Soundcloud playlist api")
	limit  = flag.String("limit", "10", "Limit items")
	offset = flag.String("offset", "0", "Offset")
)

type Playlist struct {
	Id           int     `json:"id"`
	Title        string  `json:"title"`
	ArtworkUrl   string  `json:"artwork_url"`
	Description  string  `json:"description"`
	LastModified string  `json:"last_modified"`
	Tracks       []Track `json:"tracks"`
}

type Track struct {
	Id             int    `json:"id"`
	Title          string `json:"title"`
	ArtworkUrl     string `json:"artwork_url"`
	AttachmentsUri string `json:"attachments_uri"`
	CreatedAt      string `json:"created_at"`
	Description    string `json:"description"`
	Duration       int    `json:"Duration"`
	PermalinkUrl   string `json:"permalink_url"`
	StreamUrl      string `json:"stream_url"`
	Uri            string `json:"uri"`
	VideoUrl       string `json:"video_url"`
	Stream         string
}

type Stream struct {
	Mp3Url        string `json:"http_mp3_128_url"`
	PreviewMp3Url string `json:"preview_mp3_128_url"`
}

func (t *Track) setStream() {
	t.StreamUrl = t.StreamUrl + "&client_id=" + *clientId
	ti, err := time.Parse("2006/01/02 15:04:05 -0700", t.CreatedAt)
	if err != nil {
		t.CreatedAt = time.RFC822
	} else {
		t.CreatedAt = ti.Format(time.RFC822)
	}
	return
	/*
		if t.StreamUrl != "" {
			var stream Stream
			r, err := http.Get(t.StreamUrl + "&client_id=" + *clientId)
			fmt.Println(t.StreamUrl)
			panic(t.StreamUrl)
			if err != nil {
				panic(err)
			}

			defer r.Body.Close()
			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				panic(err)
			}

			if err := json.Unmarshal(body, &stream); err != nil {
				panic(err)
			}
			t.Stream = stream
		}
	*/
}

func getPlaylist() Playlist {
	var playlist Playlist
	r, err := http.Get(*url + "&limit=" + *limit + "&offset=" + *offset)

	if err != nil {
		panic(err)
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &playlist); err != nil {
		panic(err)
	}

	return playlist
}

// Set stream on every track.
func (p *Playlist) preparePlaylist() {
	for k, _ := range p.Tracks {
		p.Tracks[k].setStream()
	}
}

func main() {
	flag.Parse()
	urlValidated, err := regexp.MatchString("https://api.soundcloud.com/playlists/.+", *url)
	if err != nil {
		panic(err)
	}

	if urlValidated == false {
		fmt.Println("Not an valid Soundcloud url, needs to be https://api.soundcloud.com/playlists/{id}..")
		return
	}

	playlist := getPlaylist()
	playlist.preparePlaylist()

	t := template.Must(template.ParseFiles("template.xml"))
	err = t.Execute(os.Stdout, playlist)

	if err != nil {
		fmt.Println("Template err", err)
	}
}
