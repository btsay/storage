package crawl

import (
	"time"

	"github.com/btlike/storage/parser"
	"github.com/btlike/storage/utils"
)

type torrentSearch struct {
	Name       string
	Length     int64
	CreateTime time.Time
}

func createElasticIndex(metaInfo parser.MetaInfo) (err error) {
	if metaInfo.InfoHash == "" ||
		metaInfo.Info.Name == "" ||
		metaInfo.Info.Length == 0 {
		return
	}
	data := torrentSearch{
		Name:       metaInfo.Info.Name,
		Length:     metaInfo.Info.Length,
		CreateTime: time.Now(),
	}
	_, err = utils.Config.ElasticClient.Index().
		Index("torrent").
		Type("infohash").
		Id(metaInfo.InfoHash).
		BodyJson(data).
		Refresh(false).
		Do()
	return
}
