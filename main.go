package main

import (
	httptransport "github.com/go-kit/kit/transport/http"
	"log"
	"net/http"
)

func main() {
	var searchSVC SearchService
	searchSVC = searchService{}
	var countSVC CountService
	countSVC = countService{}

	albumSearchHandler := httptransport.NewServer(
		makeAlbumSearchEndpoint(searchSVC),
		decodeAlbumSearchRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountInAlbumEndpoint(countSVC),
		decodeCountInAlbumRequest,
		encodeResponse,
	)

	http.Handle("/albumsearch", albumSearchHandler)
	http.Handle("/count", countHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
