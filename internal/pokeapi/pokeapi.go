package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pokedexcli/internal/pokecache"
)

type locationAreas struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

const LocationStartURL = "https://pokeapi.co/api/v2/location-area/"
const cachemsInterval = 5000

var cache *pokecache.Cache

func init() {
	cache = pokecache.NewCache(cachemsInterval)
}

func LocationPage(url string) (locations []string, nextUrl string, prevUrl string, err error) {
	bytes, err := getBuffer(url)
	if err != nil {
		return []string{}, "", "", err
	}
	var page locationAreas
	if err := json.Unmarshal(bytes, &page); err != nil {
		return []string{}, "", "", err
	}
	if page.Next == nil {
		nextUrl = ""
	} else {
		nextUrl = *page.Next
	}
	if page.Previous == nil {
		prevUrl = ""
	} else {
		prevUrl = *page.Previous
	}
	locations = make([]string, len(page.Results))
	for i, result := range page.Results {
		locations[i] = result.Name
	}
	return locations, nextUrl, prevUrl, nil
}

func getBuffer(url string) ([]byte, error) {
	bytes, ok := cache.Get(url)
	if ok {
		return bytes, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		return nil, fmt.Errorf("got status code %d", res.StatusCode)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	cache.Add(url, data)
	return data, nil
}
