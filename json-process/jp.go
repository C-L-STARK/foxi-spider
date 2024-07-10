package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
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

type Pkg struct {
	Pkg       string `json:"pkg"`
	Title string `json:"title"`
	LatestUpdate     string `json:"lu"`
	DownloadURLList       []DownloadURL `json:"download"`
}

type PathParam struct {
	Params Pkg `json:"params"`
}

func Map[T, U any](ts []T, f func(T) U) []U {
    us := make([]U, len(ts))
    for i := range ts {
        us[i] = f(ts[i])
    }
    return us
}

func transformToPageParam(item Item) *PathParam {
	// demo: https://foxirj.com/zettlr-win.html
	url, err := url.Parse(item.Url)
	
	ss := strings.Split(url.Path, "/")

	if err != nil {
		fmt.Println(err.Error())
		return nil;
	}
	
	pkg := Pkg {
		Pkg: strings.Split(ss[1], ".")[0],
		Title: "Free " + item.Title,
		LatestUpdate: item.LatestUpdate,
		DownloadURLList: item.DownloadURLList,
	}

	pathParam := PathParam {
		Params: pkg,
	}

	return &pathParam
}

func op(fileName string, outputFileName string) {
	b, e := os.ReadFile(fileName)
	if e != nil {
		fmt.Println(e.Error())
		return 
	}

	var itemList []Item

	if err := json.Unmarshal(b, &itemList); err != nil {
		fmt.Println(err.Error())
		return
	}
	
	// parsed object
	pathParamList := Map(itemList, transformToPageParam)
	j, err_1 := json.Marshal(pathParamList)
	if err_1 != nil {
		fmt.Println(err_1.Error())
		return
	}
	
	os.Create(outputFileName)
	
	if err := os.WriteFile(outputFileName, j, 0777); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("处理成功")
}

func generateMarkdownLink(fileName string, outputFileName string) {
	b, e := os.ReadFile(fileName)
	if e != nil {
		fmt.Println(e.Error())
		return 
	}

	var itemList []Item

	if err := json.Unmarshal(b, &itemList); err != nil {
		fmt.Println(err.Error())
		return
	}

	// 在原来的基础上直接生成就行：
	// Url       string `json:"url"`
	// Title string `json:"title"`

	os.Create(outputFileName)
	f, _ := os.OpenFile(outputFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)

	// 1.0 得到 URL
	currentColumn := 1
	for i := range itemList {
		item := itemList[i]

		url, err := url.Parse(item.Url)
	
		ss := strings.Split(url.Path, "/")

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		path := strings.Split(ss[1], ".")[0]
		var fullPath string
		if strings.Contains(fileName, "windows") {
			fullPath = "/software/windows/" + path
		} else {
			fullPath = "/software/macos/" + path
		}

		var links string
		if currentColumn == 6 {
			links = fmt.Sprintf("| [%s](%s){:target=\"_blank\"} |\n", item.Title, fullPath)
			currentColumn = 1
		} else {
			links = fmt.Sprintf("| [%s](%s){:target=\"_blank\"} ", item.Title, fullPath)
			currentColumn ++
		}

		_, _ = f.WriteString(links)
	}
	
	fmt.Println("处理成功")
}

func generateMenu(fileName string, outputFileName string) {
	b, e := os.ReadFile(fileName)
	if e != nil {
		fmt.Println(e.Error())
		return 
	}

	var itemList []Item

	if err := json.Unmarshal(b, &itemList); err != nil {
		fmt.Println(err.Error())
		return
	}

	os.Create(outputFileName)
	f, _ := os.OpenFile(outputFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)

	// 1.0 得到 URL
	for i := range itemList {
		item := itemList[i]

		url, err := url.Parse(item.Url)
	
		ss := strings.Split(url.Path, "/")

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		path := strings.Split(ss[1], ".")[0]
		var fullPath string // 完整的路径，如果生成英文菜单，这里需要加上 /en/
		if strings.Contains(fileName, "windows") {
			fullPath = "/software/windows/" + path
		} else {
			fullPath = "/software/macos/" + path
		}

		links := fmt.Sprintf("{ text: '%s', link: '%s' },\n", item.Title, fullPath)
		_, _ = f.WriteString(links)
	}
	
	fmt.Println("处理成功")
}

func main() {
	// PathParam JSON 生成
	// // For windows
	// op("/Users/x.stark.dylan/Desktop/fastx-ai.com/foxi-spider/json-process/windows.json", "./windows-pkg.json")
	// // For macos
	// op("/Users/x.stark.dylan/Desktop/fastx-ai.com/foxi-spider/json-process/macos.json", "./macos-pkg.json")

	// 软件主页生成
	// generateMarkdownLink("/Users/x.stark.dylan/Desktop/fastx-ai.com/foxi-spider/json-process/windows.json", "./windows-home-list.txt")
	// generateMarkdownLink("/Users/x.stark.dylan/Desktop/fastx-ai.com/foxi-spider/json-process/macos.json", "./macos-home-list.txt")

	// 菜单生成
	generateMenu("/Users/x.stark.dylan/Desktop/fastx-ai.com/foxi-spider/json-process/windows.json", "./windows-menu-list.txt")
	generateMenu("/Users/x.stark.dylan/Desktop/fastx-ai.com/foxi-spider/json-process/macos.json", "./macos-menu-list.txt")

}