package common

import (
	"bytes"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"

	"github.com/MakeGolangGreat/telegraph-go"
)

// Save 爬取内容，保存到Telegraph，然后返回线上链接的入口函数。
func Save(url string, page *telegraph.Page) (link string, err error) {
	err2 := getContent(url, page)
	errHandler("获取内容失败", err2)
	link, err = page.CreatePage()

	return
}

// 尝试用Body中获取内容
func getContent(url string, data *telegraph.Page) error {
	spider := colly.NewCollector()
	extensions.RandomUserAgent(spider)
	extensions.Referer(spider)

	spider.OnResponse(func(res *colly.Response) {
		dom, err := goquery.NewDocumentFromReader(bytes.NewReader((res.Body)))
		errHandler("初始化goquery失败", err)

		// 标题
		dom.Find("head title").Each(func(_ int, s *goquery.Selection) {
			data.Title = s.Text()
		})

		// 通用模式没办法准确获取到作者用户名
		data.AuthorName = "MakeGolangGreat"

		// 广播内容
		dom.Find("body").Each(func(_ int, s *goquery.Selection) {
			text := s.Text()
			data.Data += text
		})

		// 通用模式无法得到文章发表时间，只记录备份时间
		data.Data += "<br/><blockquote>备份时间：" + time.Now().Local().String() + "</blockquote>"
	})

	var err error
	// Set error handler
	spider.OnError(func(r *colly.Response, wrong error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", "\nError:", wrong)
		err = wrong
	})

	spider.Visit(url)

	return err
}

func errHandler(msg string, err error) {
	if err != nil {
		fmt.Printf("%s - %s\n", msg, err)
	}
}
