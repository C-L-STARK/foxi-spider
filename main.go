package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gocolly/colly"
)


type DownloadURL struct {
	Type 	string `json:"type"`
	Url       string `json:"url"`
	Password string `json:"pwd"`
}

type Item struct {
	Url       string `json:"url"`
	Title string `json:"title"`
	LatestUpdate     string `json:"lu"`
	DownloadURLList       []DownloadURL `json:"download"`
}

func main() {
	// 存放储存结果
	itemList := make([]Item, 0)

	c := colly.NewCollector()

	// 子界面获取器
	detailCollector := c.Clone()
	// 获得子界面信息
	detailCollector.OnHTML(".post-content", func (e *colly.HTMLElement)  {
		url := e.Request.URL.String()
		title := e.ChildText(".post-title")
		content := e.ChildText(".c-alert")

		downloadURLList := make([]DownloadURL, 0)
		// 循环获得内容
		e.ForEach(".c-downbtn", func(_ int, el *colly.HTMLElement) {
			passWord := el.ChildText(".c-downbtn-pwd-key")
			downloadURL := el.ChildAttr("a", "href")
			t := el.ChildAttr("img", "alt")

			// 有下载密码，才认为是完整的下载链接
			if len(passWord) > 0 {
				download := DownloadURL {
					Url: downloadURL,
					Password: passWord,
					Type: t,
				}
	
				downloadURLList = append(downloadURLList, download)	
			}
		})
		
		item := Item {
			Url: url,
			Title: title,
			LatestUpdate: content,
			DownloadURLList: downloadURLList,
		}

		fmt.Println(item.Title)

		// 增加到全局
		if len(downloadURLList) > 0 {
			// 仅有下载地址的时候才进行下载
			itemList = append(itemList, item)	
		}

		time.Sleep(3 * time.Second)
	})

	// 所有的分类页一级列表 
	c.OnHTML("a[cp-post-thumbnail-a]", func(e *colly.HTMLElement) {
		// 继续访问子页面
		detailCollector.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println("Visiting", r.URL)
	})

	// Windows 软件爬取
	// c.Visit("https://foxirj.com/category/windows")

	// for i := 2; i <= 85; i++ {
	// 	time.Sleep(1 * time.Second)
	// 	c.Visit("https://foxirj.com/category/windows/page/" + strconv.Itoa(i))
	// }

	// Macos 软件爬取
	c.Visit("https://foxirj.com/category/macos")

	for i := 2; i <= 100; i++ {
		time.Sleep(1 * time.Second)
		c.Visit("https://foxirj.com/category/macos/page/" + strconv.Itoa(i))
	}

	// 访问结束，生成 JSON 内容
	b, e := json.Marshal(itemList)
	if e != nil {
		fmt.Print(e.Error())
		return 
	}

	os.Create("./macos.json")
	err := os.WriteFile("./macos.json", b, 0777)
	
	if err != nil {
		fmt.Printf("写入异常 %v", err)
   	}
	
}
