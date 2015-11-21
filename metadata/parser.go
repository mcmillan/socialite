package metadata

import (
	"errors"
	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func ParseURL(url string) (md Metadata, err error) {
	logFields := log.Fields{
		"package": "metadata",
		"url":     url,
	}

	log.WithFields(logFields).Info("Looking for metadata...")

	html, realURL, err := fetchPageHTML(url)

	if err != nil {
		log.WithFields(logFields).Error(err)
		return
	}

	title, err := findTitle(html)

	if err != nil {
		log.WithFields(logFields).Error(err)
		return
	}

	md = Metadata{
		RealURL: realURL,
		Title:   title,
	}

	return
}

func fetchPageHTML(url string) (io.Reader, string, error) {
	res, err := httpClient.Get(url)

	if err != nil {
		return nil, "", err
	}

	return res.Body, res.Request.URL.String(), err
}

func findTitle(r io.Reader) (title string, err error) {
	doc, err := html.Parse(r)

	if err != nil {
		return
	}

	title = findTwitterTitle(doc)
	if title != "" {
		return
	}

	title = findOpenGraphTitle(doc)
	if title != "" {
		return
	}

	title = findHTMLTitle(doc)
	if title != "" {
		return
	}

	err = errors.New("Unable to ascertain title")

	return
}

func findOpenGraphTitle(doc *html.Node) string {
	el, found := scrape.Find(doc, func(n *html.Node) bool {
		if n.DataAtom == atom.Meta {
			return scrape.Attr(n, "property") == "og:title" && scrape.Attr(n, "content") != ""
		}

		return false
	})

	if !found {
		return ""
	}

	return scrape.Attr(el, "content")
}

func findTwitterTitle(doc *html.Node) string {
	el, found := scrape.Find(doc, func(n *html.Node) bool {
		if n.DataAtom == atom.Meta {
			return scrape.Attr(n, "name") == "twitter:title" && scrape.Attr(n, "content") != ""
		}

		return false
	})

	if !found {
		return ""
	}

	return scrape.Attr(el, "content")
}

func findHTMLTitle(doc *html.Node) string {
	el, found := scrape.Find(doc, scrape.ByTag(atom.Title))

	if !found {
		return ""
	}

	return scrape.Text(el)
}
