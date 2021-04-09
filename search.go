package gomdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// SearchResult contains partial information about a Film, as returned by Search.
type SearchResult struct {
	Title  string
	Year   int
	ID     string
	Type   string
	Poster string
}

// SearchOpts contains search options, such as filtering by year of release. It can be passed to Search.
type SearchOpts struct {
	Year int
	Type string
}

// Search returns an array of results for a given title. It returns an error if no film could be found.
func Search(query string, opts *SearchOpts) ([]SearchResult, error) {
	page, err := search(query, opts, 1)
	if err != nil {
		return nil, err
	}

	// TODO: retrieve further pages

	var result []SearchResult
	for _, item := range page.Search {
		film := SearchResult{
			Title:  item.Title,
			ID:     item.ID,
			Type:   item.Type,
			Poster: item.Poster,
		}
		y, err := strconv.Atoi(item.Year)
		if err != nil {
			// TODO: handle parse error (e.g. YEAR1-YEAR2)
			y = 3000
		}
		film.Year = y
		result = append(result, film)
	}
	return result, nil
}

type searchPage struct {
	Response string
	Error    string
	Total    int `json:"totalResults,string"`
	Search   []struct {
		Title  string
		Year   string
		ID     string `json:"imdbID"`
		Type   string
		Poster string
	}
}

func search(query string, opts *SearchOpts, page int) (searchPage, error) {
	url := fmt.Sprintf(searchURLFormat, APIKey, url.QueryEscape(query))
	if opts != nil {
		if opts.Year != 0 {
			url += "&y=" + strconv.Itoa(opts.Year)
		}
		if opts.Type != "" {
			// TODO: escape
			url += "&type=" + opts.Type
		}
	}

	resp, err := http.Get(url + "&page=" + strconv.Itoa(page))
	if err != nil {
		return searchPage{}, err
	}
	defer resp.Body.Close()

	var result searchPage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return searchPage{}, fmt.Errorf("could not decode response: %v", err)
	}
	if result.Response != "True" {
		return searchPage{}, errors.New(result.Error)
	}
	return result, nil
}
