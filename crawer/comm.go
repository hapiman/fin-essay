package crawer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hapiman/fin-essay/utils"
)

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
				if len(es.Time) <= 10 {
					timeStr = fmt.Sprintf("%s 00:00:01", es.Time)
				}
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
					if len(essay.Time) == 10 {
						cells[3] = fmt.Sprintf("%s 00:00:01", essay.Time)
					} else {
						cells[3] = essay.Time
					}
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

// 启动任务机器，获取最新咨询
func StartTaskRobot() {
	for {

		fmt.Println("get wdzj information")
		essaysWdzj := Grab_WDZJ()
		WriteEssay(essaysWdzj, "wdzj")

		// fmt.Println("get iyiou information")
		essaysIyiou := GrabIyiou()
		WriteEssay(essaysIyiou, "iyiou")

		time.Sleep(time.Minute * 30)
	}
}
