package crawer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hapiman/fin-essay/utils"
	"github.com/tidwall/gjson"
)

func parseHuXiu(content string) []Essay {
	essayList := []Essay{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".mod-b").Each(func(i int, s *goquery.Selection) {
		tEle := s.Find(".mod-thumb a")
		title, _ := tEle.Attr("title")
		href, _ := tEle.Attr("href")
		author := strings.TrimSpace(s.Find(".mob-ctt .mob-author .author-name").Text())
		time := strings.TrimSpace(s.Find(".mob-ctt .mob-author .time").Text())
		// 构造访问地址
		href = fmt.Sprintf("%s%s", utils.URLHuXiuBasic, href)
		essayList = append(essayList, Essay{
			Title: title, Url: href, Author: author, Time: time,
		})
	})

	return essayList
}

func grabHuXiu(pageNo int) string {
	const huxiuApi = utils.URLHuXiuApi
	client := &http.Client{}
	params := fmt.Sprintf("huxiu_hash_code=217a7d7fac22ae0af7c7c8cb11968556&page=%d&catId=102&last_dateline=", pageNo)
	req, err := http.NewRequest("POST", huxiuApi, strings.NewReader(params))
	if err != nil {
		// handle error
		fmt.Println("TestHuXiu err =>", err.Error())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err :", err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	value := gjson.Get(string(body), "data.data").Str
	return value
}

func grabHuXiuThisWeek() []Essay {
	pageNo := 1
	essays := []Essay{}
	loop := true
	for loop {
		content := grabHuXiu(pageNo)
		temps := parseHuXiu(content)
		for _, item := range temps {
			if strings.Index(item.Time, "前") == -1 {
				loop = false
				break
			}
			essays = append(essays, item)
		}

		pageNo++
	}
	return essays
}

func Grab_HuXiu() []Essay {
	res := []Essay{}
	res = grabHuXiuThisWeek()
	return res
}
