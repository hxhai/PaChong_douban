package spider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

//json转struct并筛选需要的字段
type LiuLangDiQiu struct {
	Code int64 `json:"code"`
	Data struct {
		Replies []struct {
			Content struct {
				Message string `json:"message"` //一级评论
				Plat    int64  `json:"plat"`
			} `json:"content"`
			Replies []struct {
				Content struct {
					Message string `json:"message"` //二级评论
					Plat    int64  `json:"plat"`
				} `json:"content"`
			} `json:"replies"`
		} `json:"replies"`
	} `json:"data"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
}

func Liulang() {
	//最新评论url
	var url string = "https://api.bilibili.com/x/v2/reply/main?&next=0&type=1&oid=370383499&mode=2"
	//最热评论url
	// var url1 string = "https://api.bilibili.com/x/v2/reply/main?&next=0&type=1&oid=370383499&mode=3"
	//发送请求
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("请求错误", err)
	}
	req.Header.Add("cookie", "nostalgia_conf=-1; _uuid=9AD9DF410-B126-8214-3BFE-A105916C10993B35841infoc; buvid3=7A01EE4E-2981-D955-542E-7AC123573E1735596infoc; b_nut=1656930336; buvid4=22CD8092-15C7-B1C9-7920-C241BDF5632D35596-022070418-XM+G1Gay5PBE7QOs5aMM5w%3D%3D; CURRENT_BLACKGAP=0; blackside_state=0; rpdid=|(k||RlY)|YR0J'uYlYmJ|Jlm; fingerprint=b16acaf7eb1d9b62b633692519dc3cb9; buvid_fp_plain=undefined; DedeUserID=125869224; DedeUserID__ckMd5=7115236cffb94247; buvid_fp=b16acaf7eb1d9b62b633692519dc3cb9; CURRENT_QUALITY=112; i-wanna-go-back=-1; b_ut=5; CURRENT_FNVAL=4048; bp_video_offset_125869224=699781675960238200; bsource=search_baidu; SESSDATA=7d6db56d%2C1678000966%2C7a804%2A91; bili_jct=de2d49fe0a0c5a250a8e42cd142d514b; innersign=0; b_lsid=85D710F68_1831639A553; sid=615qfidf; PVID=2")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	resp, err := client.Do(req) //http请求并返回响应
	if err != nil {
		fmt.Println("请求失败", err)
	}
	defer resp.Body.Close()
	//解析
	bodyText, err := ioutil.ReadAll(resp.Body) //读取body
	if err != nil {
		fmt.Println("io err", err)
	}
	fmt.Printf("bodyText: %v\n", bodyText)
	var resultList LiuLangDiQiu
	_ = json.Unmarshal(bodyText, &resultList) //解析json并把解析出的数据存储到resultList
	for _, result := range resultList.Data.Replies {
		fmt.Println("一级评论：", result.Content.Message)
		for _, reply := range result.Replies {
			fmt.Println("二级评论：", reply.Content.Message)
		}
	}
}
