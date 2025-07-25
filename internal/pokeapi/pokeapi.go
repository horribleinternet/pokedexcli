package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func LocationPage(url string) (locations []string, nextUrl string, prevUrl string, err error) {
	res, err := http.Get(url)
	if err != nil {
		return []string{}, "", "", err
	}
	defer res.Body.Close()
	if res.StatusCode/100 != 2 {
		return []string{}, "", "", fmt.Errorf("got status code %d", res.StatusCode)
	}
	decoder := json.NewDecoder(res.Body)
	var page locationAreas

	if err := decoder.Decode(&page); err != nil {
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
