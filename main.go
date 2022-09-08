package main

import (
	"pachong_douban/spider"
)

func main() {
	//静态数据爬取
	spider.Douban_movie()

	//动态数据爬取
	spider.Liulang()

	//并发爬取
	spider.Biji()
}
