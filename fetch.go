package gomdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// APIKey is a global variable containing the key for the OMDB service. The key must be set before using the
// package. You can get a free key here: http://omdbapi.com/apikey.aspx
var APIKey string

const (
	getURLFormat    = "http://www.omdbapi.com/?apikey=%s&i=%s"
	searchURLFormat = "http://www.omdbapi.com/?apikey=%s&s=%s"
)

// Film contains full information about a film, as returned by Fetch.
type Film struct {
	Title  string
	Year   int
	Genre  string
	Poster string
}

// Fetch returns a single film based on the ID provided. It returns an error if the film could not be found
// (e.g. because of an incorrect ID).
func Fetch(id string) (Film, error) {
	resp, err := http.Get(fmt.Sprintf(getURLFormat, APIKey, url.QueryEscape(id)))
	if err != nil {
		return Film{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Response string
		Error    string
		Title    string
		Year     string
		Genre    string
		Poster   string
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Film{}, fmt.Errorf("could not decode response: %v", err)
	}
	if result.Response == "False" {
		return Film{}, errors.New(strings.TrimSuffix(result.Error, "."))
	}

	film := Film{
		Title:  result.Title,
		Genre:  result.Genre,
		Poster: result.Poster,
	}
	y, err := strconv.Atoi(result.Year)
	if err != nil {
		return film, fmt.Errorf("could not parse year: %v", err)
	}
	film.Year = y
	return film, nil
}
