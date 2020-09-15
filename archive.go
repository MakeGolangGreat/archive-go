package archive

import (
	"errors"
	"log"
	"regexp"

	"github.com/fatih/color"

	"github.com/MakeGolangGreat/archive-go/common"
	"github.com/MakeGolangGreat/archive-go/douban"
	"github.com/MakeGolangGreat/archive-go/weibo"
	"github.com/MakeGolangGreat/archive-go/weixin"
	"github.com/MakeGolangGreat/archive-go/zhihu"

	"github.com/MakeGolangGreat/telegraph-go"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

// Save 是一个备份函数，将链接内的文本抓取然后备份到Telegraph，然后返回一个Telegraph链接。
func Save(updateText string, token string, attachInfo *telegraph.NodeElement, more *More) (msg string, err error) {
	// 手动给文本尾部增加一个空格。因为还不确定如何匹配文本中只有一个链接的字符串。
	// #HELP
	updateText += " "

	linkRegExp := regexp.MustCompile(`(http.*?)\s`)

	// 如果能匹配到某个链接
	// TODO 没有考虑到文章中有多个链接的可能，只是匹配了第一个
	if linkRegExp.MatchString(updateText) {
		replyMessage := ""
		// 拿到链接，但有可能是个错误的链接。

		matchURL := linkRegExp.FindAllSubmatch([]byte(updateText), -1)
		link := string(matchURL[0][1])

		// 如果是 telegra.ph http://archive.org/ https://archive.is/ 的链接，那么就不需要备份了。
		isArchived := regexp.MustCompile(`telegra\.ph|archive\.`)

		if isArchived.MatchString(link) {
			return "", errors.New("备份链接是telegra.ph|archive.org...链接，因此不需要备份")
		}

		page := &telegraph.Page{
			AccessToken: token,
			AuthorURL:   link,
			AuthorName:  projectName,
			AttachInfo:  attachInfo,
		}

		var err error

		if zhihu.IsZhihuLink(link) {
			color.Green("监测到知乎链接")
			replyMessage, err = zhihu.Save(link, page)
		} else if douban.IsDoubanLink(link) {
			color.Green("监测到豆瓣链接")
			replyMessage, err = douban.Save(link, page)
		} else if weibo.IsWeiboLink(link) {
			color.Green("监测到微博链接")
			replyMessage, err = weibo.Save(link, page)
		} else if weixin.IsWeixinLink(link) {
			color.Green("监测到微信链接")
			replyMessage, err = weixin.Save(link, page)
		} else {
			if !more.IncludeAll {
				// 如果不包含所有的链接，也就是不处理未适配的链接
				return "", errors.New("不处理未适配的链接")
			}

			color.Green("未适配该链接，走通用逻辑")
			replyMessage, err = common.Save(link, page)
		}

		if err != nil {
			return "", err
		}

		return replyMessage, nil
	}

	return "", errors.New("没有检测到链接")
}

// Text 是一个备份函数，将传递过来的文本备份到Telegraph，不管里面有没有链接，全部当成文本备份
// 然后返回一个Telegraph链接
func Text(updateText string, token string) (msg string, err error) {
	page := &telegraph.Page{
		AccessToken: token,
		AuthorURL:   projectLink,
		AuthorName:  projectName,
		Title:       "内容备份",
		Data:        updateText + projectDesc,
	}

	link, err := page.CreatePage()
	if err != nil {
		return "", err
	}

	return link, nil
}

// 检查此链接是否之前已经备份过，如果备份过，直接返回上次备份的链接
// 但不确定如何实现。关键在于如何保存每次的记录。本地数据库？那意味着将要长久地租一台服务器...
// 每次将保存记录保存在一个telegra.ph文章里？那么并发将是个问题，毕竟每次都要先读取telegra.ph链接来获取记录以及每次都要编辑telegra.ph文章。太频繁了。
func checkExist(link string) {
}
