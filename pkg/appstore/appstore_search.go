package appstore

import (
	"fmt"
	"github.com/majd/ipatool/pkg/http"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

type SearchResult struct {
	Count   int   `json:"resultCount,omitempty"`
	Results []App `json:"results,omitempty"`
}

func (a *appstore) Search(term, countryCode, deviceFamily string, limit int64) error {
	if StoreFronts[countryCode] == "" {
		return ErrorInvalidCountryCode
	}

	request, err := a.searchRequest(term, countryCode, deviceFamily, limit)
	if err != nil {
		return errors.Wrap(err, ErrorCreateRequest.Error())
	}

	res, err := a.searchClient.Send(request)
	if err != nil {
		return errors.Wrap(err, ErrorRequest.Error())
	}

	if res.StatusCode != 200 {
		a.logger.Debug().Interface("data", res.Data).Int("status", res.StatusCode).Send()
		return ErrorRequest
	}

	a.logger.Info().Int("count", res.Data.Count).Array("apps", Apps(res.Data.Results)).Send()
	return nil
}

func (a *appstore) searchRequest(term, countryCode, deviceFamily string, limit int64) (http.Request, error) {
	searchURL, err := a.searchURL(term, countryCode, deviceFamily, limit)
	if err != nil {
		return http.Request{}, errors.Wrap(err, ErrorURL.Error())
	}

	return http.Request{
		URL:            searchURL,
		Method:         http.MethodGET,
		ResponseFormat: http.ResponseFormatJSON,
	}, nil
}

func (a *appstore) searchURL(term, countryCode, deviceFamily string, limit int64) (string, error) {
	var entity string

	switch deviceFamily {
	case DeviceFamilyPhone:
		entity = "software"
	case DeviceFamilyPad:
		entity = "iPadSoftware"
	default:
		return "", ErrorInvalidDeviceFamily
	}

	params := url.Values{}
	params.Add("entity", entity)
	params.Add("limit", strconv.Itoa(int(limit)))
	params.Add("media", "software")
	params.Add("term", term)
	params.Add("country", countryCode)

	return fmt.Sprintf("https://%s%s?%s", iTunesAPIDomain, iTunesAPIPathSearch, params.Encode()), nil
}