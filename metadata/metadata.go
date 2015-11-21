package metadata

import (
	"net/http"
	"net/http/cookiejar"
	"time"

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
