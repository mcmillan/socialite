package metadata

import (
	"errors"
	"io"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

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
