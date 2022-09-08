package spider

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
)

const (
	USERNAME = "root"         //用户名
	PASSWORD = "root"         //密码
	HOST     = "127.0.0.1"    //主机
	PORT     = "3306"         //端口
	DBNAME   = "douban_movie" //数据库名
)

var DB *sql.DB
var t int = 1

type MovieData struct {
	Title      string `json:"title"`      //电影名
	Director   string `json:"director"`   //导演
	Picture    string `json:"picture"`    //封面
	Actor      string `json:"actor"`      //主演
	Year       string `json:"year"`       //年份
	Score      string `json:"score"`      //评分
	Evaluation string `json:"evaluation"` //评价
	Quote      string `json:"quote"`      //引言
}

//连接数据库
func InitDB() {
	//拼接数据源root:root@tcp(127.0.0.1:3306)/douban_movie?charset=utf8
	path := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", HOST, ":", PORT, ")/", DBNAME, "?charset=utf8"}, "")
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(10) //设置连接可以重用的最长时间
	DB.SetMaxIdleConns(5)     //设置最大连接数
	if err := DB.Ping(); err != nil {
		fmt.Println("数据库连接错误！")
		return
	}
	fmt.Println("数据库连接成功！")
}

//插入数据
func InsertData(movieData MovieData) bool {
	tx, err := DB.Begin() //启动事务
	if err != nil {
		fmt.Println("begin错误", err)
		return false
	}
	// 创建一个预准备语句，以便在事务中使用
	stmt, err := tx.Prepare("INSERT INTO movie_data (`Title`,`Director`,`Picture`,`Actor`,`Year`,`Score`,`Evaluation`,`Quote`) VALUES (?, ?, ?,?,?,?,?,?)")
	if err != nil {
		fmt.Println("Prepare 错误", err)
		return false
	}
	//使用给定的参数执行预准备语句
	_, err = stmt.Exec(movieData.Title, movieData.Director, movieData.Picture, movieData.Actor, movieData.Year, movieData.Score, movieData.Evaluation, movieData.Quote)
	if err != nil {
		fmt.Println("Exec fail", err)
		return false
	}
	_ = tx.Commit() //提交事务
	return true
}

func Spider(page string) {
	//1.发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://movie.douban.com/top250?start="+page, nil)
	if err != nil {
		fmt.Println("请求错误！", err)
	}
	//防止浏览器检测爬虫访问，所以加一些请求头伪造成浏览器访问
	//添加cookie和设置用户代理
	req.Header.Add("Cookie", "viewed='26832468'; bid=AQpVCG5B90w; gr_user_id=6d5152d8-6258-4a06-a6a1-db1e65dc3268; __utmc=30149280; _pk_ref.100001.4cf6=%5B%22%22%2C%22%22%2C1662368460%2C%22https%3A%2F%2Fwww.baidu.com%2Flink%3Furl%3De9xZgP2DndHRZb9GuZRKE19GW2kxYvDMLhIc5QwxFo0ALupLckZ4h5ulxkWRGI8G%26wd%3D%26eqid%3D8da4e8ca0000f48a000000066315bac6%22%5D; _pk_ses.100001.4cf6=*; ap_v=0,6.0; __utma=30149280.1527893980.1662294731.1662294731.1662368460.2; __utmb=30149280.0.10.1662368460; __utmz=30149280.1662368460.2.2.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; __utma=223695111.220523557.1662368460.1662368460.1662368460.1; __utmb=223695111.0.10.1662368460; __utmc=223695111; __utmz=223695111.1662368460.1.1.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; _pk_id.100001.4cf6=2b873ac8b1fc8973.1662368460.1.1662368959.1662368460")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	//发送 HTTP 请求并返回 HTTP 响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("请求失败: %v\n", err)
	}
	defer resp.Body.Close()
	//2.解析网页
	docDetail, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("解析失败: %v\n", err)
	}
	//3.获取节点信息
	docDetail.Find("#content > div > div.article > ol > li").
		Each(func(i int, s *goquery.Selection) { //迭代，在列表中继续找
			var data MovieData
			title := s.Find("div > div.info > div.hd > a > span:nth-child(1)").Text()
			img := s.Find("div > div.pic > a > img")
			imgTmp, ok := img.Attr("src") //获取指定属性的值
			info := s.Find("div > div.info > div.bd > p:nth-child(1)").Text()
			score := s.Find("div > div.info > div.bd > div > span.rating_num").Text()
			evaluation := s.Find("div > div.info > div.bd > div > span:nth-child(4)").Text()
			quote := s.Find("div > div.info > div.bd > p.quote > span").Text()
			if ok {
				dirextor, actor, year := InfoSpite(info)
				data.Title = title
				data.Director = dirextor
				data.Picture = imgTmp
				data.Actor = actor
				data.Year = year
				data.Score = score
				data.Evaluation = evaluation
				data.Quote = quote
				fmt.Printf("\n---第 %d 部\n %s", t, data)
				t++
				//4.保存信息
				if InsertData(data) {
					fmt.Println("————插入成功")
				} else {
					fmt.Println("插入失败")
					return
				}
			}
		})
}

//正则匹配，把导演、主演和电影年份筛选出来
func InfoSpite(info string) (dirextor, actor, year string) {
	dirextorRe, _ := regexp.Compile(`导演(.*)主`)
	dirextor = string(dirextorRe.Find([]byte(info)))
	actorRe, _ := regexp.Compile(`主演(.*)`)
	actor = string(actorRe.Find([]byte(info)))
	yearRe, _ := regexp.Compile(`(\d+)`)
	year = string(yearRe.Find([]byte(info)))
	return
}

func Douban_movie() {
	InitDB() //连接数据库
	page := 1
	for i := 0; i < 10; i++ {
		fmt.Printf("\n————正在爬取 %d 页信息————\n", page)
		Spider(strconv.Itoa(i * 25))
		page++
	}
}
