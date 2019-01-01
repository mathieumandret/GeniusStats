package main

import (
	"github.com/mathieumandret/AlbumSearchService/albumsearch"
	"github.com/mathieumandret/AlbumSearchService/lyricscount"
)

type SearchService interface {
	SearchAlbums(query string) ([]albumsearch.Album, error)
}

type searchService struct{}

// Search albums with the query string in their name in the first page of Genius
func (searchService) SearchAlbums(query string) (albums []albumsearch.Album, err error) {
	return albumsearch.SearchAlbums(query)
}

type CountService interface {
	CountInAlbum(albumURL string, term string) (int, error)
}

type countService struct{}

func (countService) CountInAlbum(albumURL string, term string) (count int, err error) {
	return lyricscount.CountInAlbum(albumURL, term)
}
