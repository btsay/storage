package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//define address
const (
	Xunlei   = "http://bt.box.n0808.com/%s/%s/%s.torrent"
	Torcache = "https://torcache.net/torrent/%s.torrent"
)

//define errors
var (
	ErrNotFound = errors.New("not found")
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

func pretty(v interface{}) {
	b, _ := json.MarshalIndent(v, " ", "  ")
	fmt.Println(string(b))
}
