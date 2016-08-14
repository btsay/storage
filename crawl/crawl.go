package crawl

import (
	"net/http"
	"net/url"
	"time"

	"github.com/btlike/storage/utils"
)

type crawl struct {
	xunleiClient *http.Client
}

func newCrawl() (c crawl, err error) {
	if utils.Config.Proxy.Enable {
		var proxyURL *url.URL
		proxyURL, err = url.Parse(utils.Config.Proxy.Address)
		if err != nil {
			return
		}
		c.xunleiClient = &http.Client{
			Timeout:   10 * time.Second,
			Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)},
		}
	} else {
		c.xunleiClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}
	return
}
