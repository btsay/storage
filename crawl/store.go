package crawl

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/btlike/database/torrent"
	"github.com/btlike/storage/parser"
	"github.com/btlike/storage/utils"
)

type torrentData struct {
	Infohash   string
	Name       string
	CreateTime time.Time
	Length     int64
	FileCount  int64

	Files []file
}

type files []file

type file struct {
	Name   string
	Length int64
}

func (a files) Len() int           { return len(a) }
func (a files) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a files) Less(i, j int) bool { return a[i].Length > a[j].Length }

//Store data in database
func Store(data parser.MetaInfo) (err error) {
	var t torrentData
	t.Infohash = data.InfoHash
	if len(t.Infohash) != 40 {
		return fmt.Errorf("store infohash len is not 40")
	}
	if data.Info.Name == "" {
		// fmt.Println("store name len is 0")
		return fmt.Errorf("store name len is 0")
	}
	t.Name = data.Info.Name
	t.CreateTime = time.Now()
	if len(data.Info.Files) == 0 {
		t.Length = data.Info.Length
		t.FileCount = 1
		t.Files = append(t.Files, file{Name: t.Name, Length: t.Length})
	} else {
		var tmpFiles files
		if len(data.Info.Files) > 5 {
			for _, v := range data.Info.Files {
				if len(v.Path) > 0 {
					t.Length += v.Length
					t.FileCount++
					tmpFiles = append(tmpFiles, file{
						Name:   v.Path[0],
						Length: v.Length,
					})
				}
			}
			sort.Sort(tmpFiles)
			if len(tmpFiles) >= 5 {
				t.Files = append(t.Files, tmpFiles[:5]...)
			} else {
				t.Files = append(t.Files, tmpFiles[:len(tmpFiles)]...)
			}
		} else {
			for _, v := range data.Info.Files {
				if len(v.Path) > 0 {
					t.Length += v.Length
					t.FileCount++
					t.Files = append(t.Files, file{
						Name:   v.Path[0],
						Length: v.Length,
					})
				}
			}
		}
	}

	b, _ := json.Marshal(t)
	if len(string(b)) > 1024 {
		return
	}
	err = insertData(t.Infohash, string(b))
	return
}

func insertData(hash string, content string) (err error) {
	switch hash[0] {
	case '0':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash0{Infohash: hash, Data: content})
	case '1':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash1{Infohash: hash, Data: content})
	case '2':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash2{Infohash: hash, Data: content})
	case '3':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash3{Infohash: hash, Data: content})
	case '4':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash4{Infohash: hash, Data: content})
	case '5':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash5{Infohash: hash, Data: content})
	case '6':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash6{Infohash: hash, Data: content})
	case '7':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash7{Infohash: hash, Data: content})
	case '8':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash8{Infohash: hash, Data: content})
	case '9':
		_, err = utils.Config.Engine.Insert(&torrent.Infohash9{Infohash: hash, Data: content})
	case 'A':
		_, err = utils.Config.Engine.Insert(&torrent.Infohasha{Infohash: hash, Data: content})
	case 'B':
		_, err = utils.Config.Engine.Insert(&torrent.Infohashb{Infohash: hash, Data: content})
	case 'C':
		_, err = utils.Config.Engine.Insert(&torrent.Infohashc{Infohash: hash, Data: content})
	case 'D':
		_, err = utils.Config.Engine.Insert(&torrent.Infohashd{Infohash: hash, Data: content})
	case 'E':
		_, err = utils.Config.Engine.Insert(&torrent.Infohashe{Infohash: hash, Data: content})
	case 'F':
		_, err = utils.Config.Engine.Insert(&torrent.Infohashf{Infohash: hash, Data: content})
	}
	return
}
