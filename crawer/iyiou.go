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

func ParseHtmlContent(content string) []Essay {
	essayList := []Essay{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".industryPostList .newestArticleList li").Each(func(i int, s *goquery.Selection) {
		aEle := s.Find("li .text a")
		boxEle := s.Find("li .text .box-lables")
		title := strings.TrimSpace(aEle.Text())
		time := strings.TrimSpace(boxEle.Find(".fr .time").Text())
		author := strings.TrimSpace(boxEle.Find(".fl .name").Text())
		href, _ := aEle.Attr("href")
		href = strings.TrimSpace(href)
		essayList = append(essayList, Essay{Title: title, Url: href, Time: time, Author: author})
	})
	return essayList
}

func grabRobot(ctx context.Context, pageNo int) string {
	log.Println("grabRobot pageNo:", pageNo)
	url := fmt.Sprintf("%s/%d.html", utils.URLIyiou, pageNo)
	content := ""
	tryTimes := 1
	for tryTimes < 5 {
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(`.industryPostList`),
			chromedp.OuterHTML(`.industryPostList`, &content),
		)
		if err == nil {
			break
		}
		log.Println(fmt.Sprintf("Try %d times", tryTimes))
		tryTimes++
	}
	return content
}

func grabThisWeek(ctx context.Context) []Essay {
	var pageNo int = 1
	essayList := []Essay{}
	for {
		temps := ParseHtmlContent(grabRobot(ctx, pageNo))
		log.Println("temps =>", temps)
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
				createdTime := fmt.Sprintf("%s 00:00:01", ele.Time)
				create_Time, _ := time.ParseInLocation("2006-01-02 15:04:05", createdTime, time.Now().Location())
				sub := time.Now().Sub(create_Time).Hours()
				if sub > 7*24 {
					break
				}
				essayList = append(essayList, ele)
			}
		}
		pageNo++
	}
	return essayList
}

func Grab() []Essay {
	// 初始化
	// 打开网页
	// 解析网页内容
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	log.Println("grab start")
	res := grabThisWeek(ctx)
	log.Println("res =>", res)
	return res
}
