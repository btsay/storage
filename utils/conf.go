package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/btlike/repository"
	"gopkg.in/olivere/elastic.v3"
)

//Config define config
var (
	Config        config
	Log           *log.Logger
	ElasticClient *elastic.Client
	Repository    repository.Repository
)

type config struct {
	Database string `json:"database"`
	Elastic  string `json:"elastic"`
	Proxy    struct {
		Enable  bool   `json:"enable"`
		Address string `json:"address"`
	} `json:"proxy"`
}

//Init utilsl
func Init() {
	initLog()
	initConfig()
	initDatabase()
	initElastic()
}

func initElastic() {
	client, err := elastic.NewClient(elastic.SetURL(Config.Elastic))
	exit(err)
	ElasticClient = client
	ElasticClient.CreateIndex("torrent").Do()
}

func initConfig() {
	f, err := os.Open("config/storage.conf")
	exit(err)
	b, err := ioutil.ReadAll(f)
	exit(err)
	var c config
	err = json.Unmarshal(b, &c)
	exit(err)
}

func initLog() {
	Log = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func initDatabase() {
	repo, err := repository.NewMysqlRepository(Config.Database, 1024, 1024)
	exit(err)
	Repository = repo
}

func exit(err error) {
	if err != nil {
		Log.Fatalln(err)
	}
}
