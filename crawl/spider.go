package crawl

import (
	"time"

	"github.com/btsay/storage/parser"
	"github.com/btsay/storage/utils"
)

//Run the spider
func Run() {
	Manager.run()
	Crawl()
}

//Crawl from xunlei ...
func Crawl() {
	worker := func(jobs <-chan string, resultChan chan<- string) {
		crawl, err := newCrawl()
		if err != nil {
			utils.Log.Printf("设置了代理，但代理地址错误：%v\n", err)
			return
		}
		for infohash := range jobs {
			if !Manager.crawStatus[Xunlei].pauseCrawl {
				//至少有一个引擎在服务时，直接删除即可，防止引擎都不服务时，疯狂删数据
				resultChan <- infohash
			}

			var data parser.MetaInfo
			var err error

			if !Manager.crawStatus[Xunlei].pauseCrawl {
				data, err = parser.DownloadTorrent(infohash, crawl.xunleiClient)
				if err != nil {
					if err == parser.ErrNotFound {
						Manager.crawStatus[Xunlei].notFoundCount++
					} else {
						Manager.crawStatus[Xunlei].refuseCount++
					}
					continue
				} else {
					//没报错，进入存储流程
					goto store
				}
			}

			if Manager.crawStatus[Xunlei].pauseCrawl {
				//预防引擎都没有服务时，直接进入下一循环
				continue
			}

		store:
			err = Store(data)
			if err != nil {
				resultChan <- infohash
				continue
			}

			//全文索引
			err = createElasticIndex(data)
			if err != nil {
				utils.Log.Println(err)
			}

			resultChan <- infohash
			Manager.storeCount++
		}
	}

	jobChan := make(chan string, DownloadChanLength)
	resultChan := make(chan string, DownloadChanLength)
	defer close(resultChan)
	for i := 0; i < 100; i++ {
		go worker(jobChan, resultChan)
	}

	go func() {
		var infohashs []string
		for infohash := range resultChan {
			if len(infohash) == 40 {
				infohashs = append(infohashs, infohash)
				if len(infohashs) >= 100 {
					err := utils.Repository.BatchDeleteInfohash(infohashs)
					infohashs = make([]string, 0)
					if err != nil {
						utils.Log.Println("delete error", err)
					}
				}
			}
		}
	}()

	for {
		if Manager.crawStatus[Xunlei].pauseCrawl {
			utils.Log.Println("全部引擎拒绝服务,暂停抓取,等待10分钟")
			time.Sleep(time.Minute * 10)
			Manager.crawStatus[Xunlei] = &crawStatus{}
		}

		pres, err := utils.Repository.BatchGetInfohash(1000)
		if err != nil {
			utils.Log.Println(err)
			time.Sleep(time.Second * 10)
			continue
		}

		if len(pres) == 0 {
			time.Sleep(time.Second * 60)
		}
		for _, v := range pres {
			jobChan <- v
		}
	}
}
