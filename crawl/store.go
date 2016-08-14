package crawl

import (
	"fmt"
	"sort"
	"time"

	"github.com/btlike/repository"
	"github.com/btlike/storage/parser"
	"github.com/btlike/storage/utils"
)

type files []repository.File

func (a files) Len() int           { return len(a) }
func (a files) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a files) Less(i, j int) bool { return a[i].Length > a[j].Length }

//Store data in database
func Store(data parser.MetaInfo) (err error) {
	var t repository.Torrent
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
		t.Files = append(t.Files, repository.File{Name: t.Name, Length: t.Length})
	} else {
		var tmpFiles files
		if len(data.Info.Files) > 5 {
			for _, v := range data.Info.Files {
				if len(v.Path) > 0 {
					t.Length += v.Length
					t.FileCount++
					tmpFiles = append(tmpFiles, repository.File{
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
					t.Files = append(t.Files, repository.File{
						Name:   v.Path[0],
						Length: v.Length,
					})
				}
			}
		}
	}
	err = utils.Repository.CreateTorrent(t)
	return
}
