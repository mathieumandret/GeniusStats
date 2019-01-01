package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	as "github.com/mathieumandret/AlbumSearchService/albumsearch"
	"net/http"
)

type CountInAlbumRequest struct {
	AlbumURL string `json:"albumURL"`
	Query    string `json:"query"`
}

type CountInAlbumResponse struct {
	Count int    `json:"count"`
	Err   string `json:"err,omitempty"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

type SearchResponse struct {
	Albums []as.Album `json:"albums"`
	Err    string     `json:"err,omitempty"`
}

func makeAlbumSearchEndpoint(svc SearchService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(SearchRequest)
		albums, err := svc.SearchAlbums(req.Query)
		if err != nil {
			return SearchResponse{albums, err.Error()}, nil
		}
		return SearchResponse{albums, ""}, nil
	}
}

func makeCountInAlbumEndpoint(svc CountService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(CountInAlbumRequest)
		count, err := svc.CountInAlbum(req.AlbumURL, req.Query)
		if err != nil {
			return CountInAlbumResponse{count, err.Error()}, nil
		}
		return CountInAlbumResponse{count, ""}, nil
	}
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeAlbumSearchRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeCountInAlbumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CountInAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}
