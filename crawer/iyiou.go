package crawer

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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
	for tryTimes < 5 { // 似乎没有起作用
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitVisible(`.industryPostList`),
			chromedp.OuterHTML(`.industryPostList`, &content),
		)
		if err == nil {
			break
		}
		fmt.Println("grab error: ", err.Error())
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
					fmt.Println("ok, go back.")
					return essayList
				}
				essayList = append(essayList, ele)
			}
		}
		pageNo++
	}
	return essayList
}

// 抓取最近7天数据
func ReadEssay(fname string) []Essay {
	// 确保文件夹和文件存在
	EnsureEssayDir()
	filePath := GetEssayFilePath(fmt.Sprintf("%s.md", fname))

	results := []Essay{} // 返回数据
	essays := utils.ReadLineEachTime(filePath)
	for _, eRow := range essays {
		eArr := strings.Split(eRow, "@@")
		if len(eArr) < 5 {
			continue
		}
		timeStr := eArr[3]
		eEntity := Essay{Title: eArr[0], Url: eArr[1], Author: eArr[2], Time: eArr[3]}
		if strings.Index(timeStr, "前") > -1 { // 比较时间
			results = append(results, eEntity)
			continue
		}
		oldT, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
		if err != nil { // handle error
			fmt.Println("time format error occurs: ", err.Error())
			continue
		}
		oldSeds := oldT.Unix()
		curSeds := time.Now().Unix()
		duSeds := int64(7 * 24 * 60 * 60)
		if curSeds-oldSeds <= duSeds {
			results = append(results, eEntity)
		}
	}
	return results
}

func appendFile(filePath string, essays []Essay) {
	// 打开文件
	f, _ := os.OpenFile(filePath, os.O_APPEND|os.O_RDONLY|os.O_WRONLY, os.ModeAppend) //打开文件
	defer f.Close()

	cStr := utils.ReadAllOnce(filePath)
	lines := strings.Split(cStr, "\n")

	for _, es := range essays {
		exsit := false
		for _, line := range lines {
			newTitle := es.Title
			cells := strings.Split(line, "@@")
			oldTitle := cells[0]
			if newTitle == oldTitle {
				exsit = true
			}
		}
		if !exsit { //创建
			timeStr := es.Time
			if strings.Index(timeStr, "前") < 0 {
				timeStr = fmt.Sprintf("%s 00:00:01", es.Time)
			}
			rowStr := fmt.Sprintf("%s@@%s@@%s@@%s@@%s\n", es.Title, es.Url, es.Author, timeStr, "xxxx")
			_, err := f.Write([]byte(rowStr))
			if err != nil {
				fmt.Println("error ->", err.Error())
			}
		}
	}
}

func WriteEssay(essays []Essay, fname string) {
	EnsureEssayDir()
	filePath := GetEssayFilePath(fmt.Sprintf("%s.md", fname))

	appendFile(filePath, essays)

	cStr := utils.ReadAllOnce(filePath)
	// fmt.Println("cStr =>", cStr)
	lines := strings.Split(cStr, "\n")

	// 判断
	for i, line := range lines {
		cells := strings.Split(line, "@@")
		if len(cells) < 5 {
			continue
		}
		oldTitle := cells[0]
		oldTime := cells[3]

		if strings.Index(oldTime, "前") < 0 {
			continue
		}
		for _, essay := range essays {
			if essay.Title == oldTitle { // 替换
				if strings.Index(essay.Time, "前") > -1 {
					cells[3] = essay.Time
				} else {
					cells[3] = fmt.Sprintf("%s 00:00:01", essay.Time)
				}
			}
		}
		lines[i] = strings.Join(cells, "@@")
	}
	output := strings.Join(lines, "\n")
	err := ioutil.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func GrabIyiou() []Essay {
	// 初始化
	// 打开网页
	// 解析网页内容
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()
	// create a timeout
	// ctx, cancel = context.WithTimeout(ctx, 90*time.Second)
	// defer cancel()
	log.Println("GrabIyiou start")
	res := grabThisWeek(ctx)
	log.Println("GrabIyiou get list =>", res)
	return res
}

// 启动任务机器
// 每隔30分钟获取一边最新咨询
func StartTaskRobot() {
	for {
		essays := GrabIyiou()
		WriteEssay(essays, "iyiou")
		time.Sleep(time.Minute * 30)
	}
}
