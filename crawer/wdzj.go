package crawer

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/hapiman/fin-essay/utils"
)

func ParseWdzjHtmlContent(content string) []Essay {
	essayList := []Essay{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".zllist li").Each(func(idx int, s *goquery.Selection) {
		aEle := s.Find("li .text h3 a")
		boxEle := s.Find("li .text .userxx .lbox span")
		aEle.Find("em").Remove()
		title := strings.TrimSpace(aEle.Text())
		author := ""
		time := ""
		boxEle.Each(func(i int, s *goquery.Selection) {
			curText := strings.TrimSpace(s.Text())
			if i == 0 {
				author = curText
			}
			childLen := boxEle.Length()
			if childLen == 3 { // 存在作者名称
				if i == 2 {
					time = curText
				}
				if i == 1 {
					author = fmt.Sprintf("%s-%s", author, curText)
				}
			} else if childLen == 1 {
				time = curText
			} else {
				if i == 1 {
					time = curText
				}
			}
		})
		href, _ := aEle.Attr("href")
		href = fmt.Sprintf("%s/%s", utils.URLWdzj, strings.TrimSpace(href))
		// xx := fmt.Sprintf("title:%s, url: %s, time: %s, author: %s", title, href, time, author)
		// fmt.Println("xx =>", xx)
		essayList = append(essayList, Essay{Title: title, Url: href, Time: time, Author: author})
	})
	return essayList
}

func grabWdzjRobot(ctx context.Context, pageNo int) string {
	log.Println("grabWdzjRobot pageNo:", pageNo)
	url := ""
	if pageNo == 1 {
		url = utils.URLWdzj
	} else {
		url = fmt.Sprintf("%s/p%d.html", utils.URLWdzj, pageNo)
	}
	content := ""
	tryTimes := 1
	for tryTimes < 5 {
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(`.listbox`),
			chromedp.OuterHTML(`.listbox`, &content),
		)
		if err == nil {
			break
		}
		log.Println(fmt.Sprintf("Try %d times", tryTimes))
		tryTimes++
	}
	return content
}

func grabWdzjThisWeek(ctx context.Context) []Essay {
	var pageNo int = 1
	essayList := []Essay{}
	for {
		temps := ParseWdzjHtmlContent(grabWdzjRobot(ctx, pageNo))
		fmt.Println("temps =>", temps)
		if len(temps) == 0 {
			break
		} else {
			for _, ele := range temps {
				// 如果存在数据 分钟前 小时前
				log.Println("ele.Time =>", ele.Time)
				ord := strings.Index(ele.Time, "分钟前")
				if ord > 0 {
					essayList = append(essayList, ele)
					continue
				}
				ord = strings.Index(ele.Time, "小时前")
				if ord > 0 {
					essayList = append(essayList, ele)
					continue
				}
				// 抓取最近7天数据
				createdTime, _ := time.ParseInLocation("2006-01-02 15:04:05", ele.Time, time.Now().Location())
				diffDays := int(time.Now().Sub(createdTime).Hours())
				fmt.Println("createdTime => ", createdTime, diffDays)
				if diffDays > 1*24 {
					return essayList
				}
				essayList = append(essayList, ele)
			}
		}
		pageNo++
	}
	return essayList
}

func Grab_WDZJ() []Essay {
	// 初始化
	// 打开网页
	// 解析网页内容
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	// create a timeout
	// ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	// defer cancel()
	log.Println("wdzj grab start")
	res := grabWdzjThisWeek(ctx)
	log.Println("res =>", res)
	return res
}
