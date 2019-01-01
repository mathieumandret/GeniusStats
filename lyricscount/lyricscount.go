package lyricscount

import (
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strings"
)

const domain = "genius.com"
const titleSelector = "div.chart_row-content"

type Track struct {
	Name   string
	URL    string
	Lyrics string
}

func sanitizeTitle(title string) string {
	return strings.TrimSuffix(title, "Lyrics")
}

func sanitizeLyrics(lyrics string) string {
	return strings.TrimSpace(strings.Replace(lyrics, "\n", " ", -1))
}

func buildCollector() colly.Collector {
	collector := colly.NewCollector(
		colly.AllowedDomains(domain),
		colly.MaxDepth(1),
	)
	return *collector
}

func fetchTrack(albumURL string) ([]Track, error) {
	var tracks []Track
	var lastTrack string
	var lyricsPerTrack = make(map[string]string)
	var err error
	collector := buildCollector()
	// When a song title is found
	collector.OnHTML(titleSelector, func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		// Set the song as lastTrack
		name := sanitizeTitle(e.ChildText("h3"))
		lastTrack = name
		// Save it in the tracks array
		tracks = append(tracks, Track{Name: name, URL: link})
		// Visit the song lyrics page
		err = collector.Visit(link)
	})
	// When lyrics are encountered, save them in the map with the last track
	// title as a key
	collector.OnHTML(".lyrics", func(e *colly.HTMLElement) {
		lyricsPerTrack[lastTrack] = sanitizeLyrics(e.Text)
	})
	collector.OnRequest(func(request *colly.Request) {
		fmt.Printf("Visiting %s \n", request.URL)
	})
	err = collector.Visit(albumURL)
	fmt.Println("Sorting lyrics")
	associateLyrics(lyricsPerTrack, tracks)
	if err != nil {
		return tracks, err
	}
	return tracks, nil
}

func associateLyrics(lyricsMap map[string]string, tracks []Track) {
	for index, track := range tracks {
		lyrics, present := lyricsMap[track.Name]
		if !present {
			continue
		}
		tracks[index].Lyrics = lyrics
	}
}

func countOccurences(term string, in string) int {
	pattern := regexp.MustCompile("(?i)" + term)
	matches := pattern.FindAllStringIndex(in, -1)
	return len(matches)
}

func allLyrics(tracks []Track) string {
	var sb strings.Builder
	for _, track := range tracks {
		sb.WriteString(track.Lyrics)
	}
	return sb.String()
}

func CountInAlbum(albumURL string, term string) (int, error) {
	tracks, err := fetchTrack(albumURL)
	fmt.Println("Grouping lyrics")
	all := allLyrics(tracks)
	fmt.Println("Counting occurences")
	return countOccurences(term, all), err
}
