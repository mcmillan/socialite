package metadata

import (
	"net/url"
)

var queryStringBlacklist = []string{
	"utm_medium",
	"utm_source",
	"utm_campaign",
	"utm_content",
}

func normalizeURL(rawURL string) (normalizedURL string) {
	url, _ := url.Parse(rawURL)

	query := url.Query()

	for _, key := range queryStringBlacklist {
		query.Del(key)
	}

	url.RawQuery = query.Encode()

	return url.String()
}
