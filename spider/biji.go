package spider

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type GormData struct {
	First_Title       string `json:"first_title"`       //一级标题
	Second_Title_Href string `json:"second_title_href"` //二级标题
	Body              string `json:"body"`              //内容
}

var Client http.Client
var wg sync.WaitGroup

func Biji() {
	url := "https://gorm.io/zh_CN/docs/"
	SpiderBiji(url, nil, 1)
	// NormalStart(url) // 单线程爬虫 672.9508ms
	// ChannelStart(url) // Channel多线程爬虫 214.6677ms
	// WaitGroupStart(url) // WaitGroup 多线程爬虫 235.2646ms
}

//单线程
func NormalStart(url string) {
	start := time.Now()
	for i := 0; i < 10; i++ {
		SpiderBiji(url, nil, i)
	}
	elapsed := time.Since(start)
	fmt.Printf("NormalStart Time %s \n", elapsed)
}

//channel多线程
func ChannelStart(url string) {
	ch := make(chan bool)
	start := time.Now()
	for i := 0; i < 10; i++ {
		go SpiderBiji(url, ch, i)
	}
	for i := 0; i < 10; i++ {
		<-ch
	}
	elapsed := time.Since(start)
	fmt.Printf("ChannelStart Time %s \n", elapsed)
}

//Waitgroup多线程
func WaitGroupStart(url string) {
	start := time.Now()
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			SpiderBiji(url, nil, i)
		}(i)
	}
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("WaitGroupStart Time %s\n ", elapsed)
}

func SpiderBiji(url string, ch chan bool, i int) {
	client := http.Client{}
	reqSpider, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	reqSpider.Header.Add("cookie", "_ga=GA1.2.1420079594.1662552623; _gid=GA1.2.1197555525.1662552623; __atuvc=1%7C36; __atuvs=63188a2e0063c413000")
	reqSpider.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	respSpider, err := client.Do(reqSpider)
	if err != nil {
		log.Fatal(err)
	}
	//2.解析网页
	docDetail, err := goquery.NewDocumentFromReader(respSpider.Body)
	if err != nil {
		fmt.Printf("解析失败: %v\n", err)
	}
	var firstTitle []string
	docDetail.Find("#sidebar > div > strong").Each(func(i int, s *goquery.Selection) {
		s2 := s.Text()
		firstTitle = append(firstTitle, s2)
	})
	fmt.Printf("一级标题: %v\n\n", firstTitle)

	//3.获取节点信息
	var secondtTitle []string
	var body []string
	docDetail.Find("#sidebar > div > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href") //获取href属性
		second_title := s.Text()
		secondtTitle = append(secondtTitle, second_title+":"+url+href) //保存内容到secondTitle切片

		//爬取内容
		rq, _ := http.NewRequest("GET", url+href, nil)
		rq.Header.Add("cookie", "_ga=GA1.2.1420079594.1662552623; _gid=GA1.2.1197555525.1662552623; __atuvc=1%7C36; __atuvs=63188a2e0063c413000")
		rq.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
		rp, _ := client.Do(rq)
		doc, _ := goquery.NewDocumentFromReader(rp.Body)
		s2 := doc.Find("#content-inner > article > div > div > div").Text()
		body = append(body, s2) //保存内容到body切片
	})
	for _, t := range secondtTitle {
		fmt.Printf("t: %v\n", t)
	}
	if ch != nil {
		ch <- true
	}
}
