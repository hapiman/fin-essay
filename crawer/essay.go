package crawer

import (
	"fmt"
	"os"
)

type Essay struct {
	Title  string `json:"title"`
	Url    string `json:"url"`
	Time   string `json:"time"`
	Author string `json:"author"`
}

func checkExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func EnsureEssayDir() {
	pwd, _ := os.Getwd()
	essayDir := fmt.Sprintf("%s%s", pwd, "/essays")
	if _, err := os.Stat(essayDir); err != nil {
		os.Mkdir(essayDir, os.ModePerm) // 创建
	}
}

func GetEssayFilePath(fname string) string {
	pwd, _ := os.Getwd()
	essayPath := fmt.Sprintf("%s/essays/%s", pwd, fname)
	if checkExists(essayPath) == false { // 创建文件
		os.Create(essayPath)
	}
	return essayPath
}
