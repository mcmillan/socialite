package metadata

import (
	"crypto/sha1"
	"io"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

type Metadata struct {
	RealURL string
	Title   string
}

func (md *Metadata) ID() string {
	h := sha1.New()
	io.WriteString(h, strings.ToLower(md.RealURL))
	return string(h.Sum(nil))
}
