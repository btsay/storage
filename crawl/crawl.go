package crawl

import (
	"net/http"
	"time"
)

type crawl struct {
	xunleiClient *http.Client
}

func newCrawl() (c crawl) {
	c.xunleiClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	return
}
