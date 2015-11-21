package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var httpClient = &http.Client{}

func request(endpoint string, method string, params map[string]string, outlet interface{}) error {
	method = strings.ToUpper(method)
	encodedParams := encodeParams(params)
	url := buildURL(method, endpoint, encodedParams)

	var body io.Reader
	if method != "GET" {
		body = strings.NewReader(encodedParams)
	}

	logFields := log.Fields{
		"url":     url,
		"method":  method,
		"package": "twitter",
	}

	log.WithFields(logFields).Info("Starting request")

	req, err := http.NewRequest(method, url, body)

	if err != nil {
		log.WithFields(logFields).Error(err)
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TWITTER_OAUTH2_BEARER_TOKEN")))

	res, err := httpClient.Do(req)

	if err != nil {
		log.WithFields(logFields).Error(err)
		return err
	}

	if res.StatusCode != 200 {
		errorBytes, _ := ioutil.ReadAll(res.Body)
		err := errors.New(string(errorBytes))
		log.WithFields(logFields).Error(err)
		return err
	}

	log.WithFields(logFields).Info("Request successful")

	return json.NewDecoder(res.Body).Decode(outlet)
}

func encodeParams(params map[string]string) string {
	data := url.Values{}

	for key, value := range params {
		if value != "" {
			data.Set(key, value)
		}
	}

	return data.Encode()
}

func buildURL(method string, endpoint string, encodedParams string) string {
	url := url.URL{
		Scheme: "https",
		Host:   "api.twitter.com",
		Path:   fmt.Sprintf("1.1/%s.json", endpoint),
	}

	if method == "GET" {
		url.RawQuery = encodedParams
	}

	return url.String()
}
