package albumsearch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/* Domain types */

type APIResponse struct {
	Meta struct {
		Status int
	}
	Response struct {
		Sections []struct {
			Type string
			Hits [] struct {
				Album APIAlbum `json:"result"`
			}
		}
	}
}

type APIAlbum struct {
	Name   string
	ArtURL string `json:"cover_art_thumbnail_url"`
	Artist APIArtist
	URL    string
}

type APIArtist struct {
	Name string
}

type Album struct {
	Name   string `json:"name"`
	ArtURL string `json:"artURL"`
	Artist string `json:"artist"`
	URL    string
}

const albumSearchURL = "https://genius.com/api/search/album?page=1&q=%s"

func buildQuery(input string) string {
	return strings.Replace(input, " ", "+", -1)
}

// Converts an album from the Genius API to a domain album model
func convertAlbum(apiAlbum APIAlbum) (album Album) {
	album.Artist = apiAlbum.Artist.Name
	album.Name = apiAlbum.Name
	album.ArtURL = apiAlbum.ArtURL
	album.URL = apiAlbum.URL
	return album
}

// Search albums with the query string in their name in the first page of Genius
func SearchAlbums(query string) (albums []Album, err error) {
	res, err := http.Get(fmt.Sprintf(albumSearchURL, buildQuery(query)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response APIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	if response.Meta.Status != 200 {
		return nil, errors.New(fmt.Sprintf("Received HTTP Code %d", response.Meta.Status))
	}
	for _, element := range response.Response.Sections {
		for _, hit := range element.Hits {
			albums = append(albums, convertAlbum(hit.Album))
		}
	}
	return albums, nil
}
