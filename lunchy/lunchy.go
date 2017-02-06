package lunchy

import (
	"encoding/json"
	"fmt"
	"github.com/rebasar/lunch/fsutil"
	"io/ioutil"
	"net/http"
)

const (
	DEFAULT_CACHE_FILENAME = ".lunchy.cache"
)

type Place struct {
	Name    string
	Uri     string
	Website string
	Aliases []string
}

func (self Place) hasAlias(s string) bool {
	for _, alias := range self.Aliases {
		if alias == s {
			return true
		}
	}
	return false
}

type Menu struct {
	ValidFrom  Date
	ValidUntil Date
	Items      []Item
}

type Item struct {
	Title       string
	Description string
	Price       uint16
}

type Client struct {
	Uri    string
	places []Place
}

func (self Client) GetPlaces() []Place {
	return self.places
}

func (self Client) GetMenu(alias string) (Menu, error) {
	result := Menu{}
	place, err := self.getPlace(alias)
	if err != nil {
		return result, err
	}
	body, err := fetch(place.Uri)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

func (self Client) getPlace(alias string) (Place, error) {
	for _, place := range self.places {
		if place.hasAlias(alias) {
			return place, nil
		}
	}
	return Place{}, fmt.Errorf("Place \"%s\" not found", alias)
}

type LunchHTTPError struct {
	StatusCode int
	Status     string
	Message    string
}

func (self LunchHTTPError) Error() string {
	if self.StatusCode == 404 {
		return self.Message
	} else {
		return fmt.Sprintf("Got status code %d from server with message: %s", self.StatusCode, self.Status)
	}
}

func fetch(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		return nil, LunchHTTPError{resp.StatusCode, resp.Status, fmt.Sprintf("No lunch found for today at %s", uri)}
	}
	return ioutil.ReadAll(resp.Body)
}

type PlaceListFetcher interface {
	fetchPlaceList() ([]byte, error)
}

type HttpPlaceListFetcher struct {
	uri string
}

func (self HttpPlaceListFetcher) fetchPlaceList() ([]byte, error) {
	return fetch(self.uri)
}

type CachePlaceListFetcher struct {
	delegate      PlaceListFetcher
	cacheFileName string
}

func (self CachePlaceListFetcher) fetchPlaceList() ([]byte, error) {
	body, err := fsutil.ReadConfigFile(self.cacheFileName)
	if err != nil {
		return self.fetchPlaceListFromDelegate()
	}
	return body, err
}

func (self CachePlaceListFetcher) fetchPlaceListFromDelegate() ([]byte, error) {
	result, err := self.delegate.fetchPlaceList()
	if err != nil {
		return result, err
	}
	err = self.updateCache(result)
	if err != nil {
		self.removeCache()
	}
	return result, err
}

func (self CachePlaceListFetcher) updateCache(data []byte) error {
	return fsutil.WriteConfigFile(data, self.cacheFileName)
}

func (self CachePlaceListFetcher) removeCache() {
	fsutil.RemoveConfigFile(self.cacheFileName)
}

func fetchPlaceList(uri string, fetcher PlaceListFetcher) ([]Place, error) {
	result := []Place{}
	body, err := fetcher.fetchPlaceList()
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(body, &result)
	return result, err
}

func NewClient(uri string) (Client, error) {
	cacheFileName, err := fsutil.GetAbsoluteFilename(DEFAULT_CACHE_FILENAME)
	if err != nil {
		return Client{}, err
	}
	return NewClientWithCache(uri, cacheFileName)
}

func NewClientWithCache(uri, cacheFileName string) (Client, error) {
	fetcher := CachePlaceListFetcher{HttpPlaceListFetcher{uri}, cacheFileName}
	places, err := fetchPlaceList(uri, fetcher)
	return Client{uri, places}, err
}
