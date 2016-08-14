## Storage
[![Build Status](https://drone.io/github.com/btlike/storage/status.png)](https://drone.io/github.com/btlike/storage/latest)

根据infohash从xunlei和torcache获取torrent metadata存储到数据库，同时增加搜索引擎全文索引



## 特性

- 分表存储，设计容量6千万到8千万，均分在16张表中
- 丢弃torrent部分字段(piece字段)，节省90%网络流量
- 引擎健康检查，全部拒绝服务时，暂停抓取
- 多线程抓取
- 支持代理



### 安装
`go get github.com/btlike/storage`

### 配置

```
{
  "database": "root:password@tcp(127.0.0.1:3306)/torrent?charset=utf8&parseTime=True&loc=Local", //数据库地址
  "elastic": "http://127.0.0.1:9200", //elasticsearch地址
  "proxy": {
    "enable": false, //是否开启代理
    "address": "http://127.0.0.1:8090" //代理地址
  }
}

```

某些情况下，ip也许会被xunlei封杀，可以采用代理模式，自建代理很简单

用这个库[github.com/elazarl/goproxy](https://github.com/elazarl/goproxy)，以下代码就可以搭建一个代理

```go
package main

import (
    "github.com/elazarl/goproxy"
    "log"
    "net/http"
)

func main() {
    proxy := goproxy.NewProxyHttpServer()
    proxy.Verbose = true
    log.Fatal(http.ListenAndServe(":8090", proxy))
}
```
