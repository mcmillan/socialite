package metadata

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/mcmillan/socialite/store"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
)

var (
	cookieJar, _ = cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
		Jar:     cookieJar,
	}
)

func ParseURL(url string) (link store.Link, err error) {
	logFields := log.Fields{
		"package": "metadata",
		"url":     url,
	}

	log.WithFields(logFields).Info("Looking for metadata...")

	html, finalURL, err := fetchPageHTML(url)

	if err != nil {
		log.WithFields(logFields).Error(err)
		return
	}

	title, err := findTitle(html)

	if err != nil {
		log.WithFields(logFields).Error(err)
		return
	}

	link = store.Link{
		Title: title,
		URL:   finalURL,
	}

	return
}

func fetchPageHTML(url string) (io.Reader, string, error) {
	res, err := httpClient.Get(url)

	if err != nil {
		return nil, "", err
	}

	finalURL := normalizeURL(res.Request.URL.String())

	return res.Body, finalURL, err
}
