package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

//define address
const (
	Xunlei   = "http://bt.box.n0808.com/%s/%s/%s.torrent"
	Torcache = "https://torcache.net/torrent/%s.torrent"
)

//define errors
var (
	ErrNotFound = errors.New("not found")
	LibUrls     = []string{
		"http://www.torrent.org.cn/Home/torrent/download.html?hash=%s",
		"http://torcache.net/torrent/%s.torrent",
		"http://torrage.com/torrent/%s.torrent",
		"http://zoink.it/torrent/%s.torrent",
		"https://178.73.198.210/torrent/%s.torrent",
		"http://d1.torrentkittycn.com/?infohash=%s",
		"http://reflektor.karmorra.info/torrent/%s.torrent",
	}
)

//DownloadXunlei torrent
func DownloadXunlei(hash string, client *http.Client) (mi MetaInfo, err error) {
	mi.InfoHash = hash
	if len(hash) != 40 {
		err = errors.New("invalid hash len")
		return
	}

	//从迅雷种子库查找
	address := fmt.Sprintf(Xunlei, hash[:2], hash[len(hash)-2:], hash)
	req0, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return
	}
	req0.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req0)
	if err != nil {
		return
	}
	if resp != nil {
		defer func() {
			// io.Copy(ioutil.Discard, resp.Body)
			resp.Body.Close()
		}()

		if resp.StatusCode == 200 {
			//解析种子
			err = mi.Parse(resp.Body)
		} else if resp.StatusCode == 404 {
			err = ErrNotFound
		} else {
			err = errors.New("refuse error")
		}
	}
	return
}

//Download torrent
func DownloadTorrent(hash string, client *http.Client) (mi MetaInfo, err error) {
	if len(hash) != 40 {
		err = errors.New("invalid hash len")
		return
	}
	mi, err = DownloadXunlei(hash, client)
	//迅雷解析成功，不用再調用後面的種子庫
	if err == nil {
		return
	}

	mi.InfoHash = hash
	//將來改用字典實現
	for _, lib_url := range LibUrls {
		address := fmt.Sprintf(lib_url, strings.ToUpper(hash))
		req0, err := http.NewRequest("GET", address, nil)
		if err != nil {
			continue
		}
		resp, err := client.Do(req0)
		if err != nil {
			continue
		}
		if resp != nil {
			defer func() {
				// io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}()

			if resp.StatusCode == 200 {
				//解析种子
				err = mi.Parse(resp.Body)
				return
			} else if resp.StatusCode == 404 {
				err = ErrNotFound
			} else {
				err = errors.New("refuse error")
			}
		}
	}
	return
}

func pretty(v interface{}) {
	b, _ := json.MarshalIndent(v, " ", "  ")
	fmt.Println(string(b))
}
