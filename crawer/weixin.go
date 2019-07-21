package crawer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hapiman/fin-essay/utils"
	"github.com/tidwall/gjson"
)

func matchUrl(wxname string) string {
	if wxname == "uni-fin" {
		return utils.URLWxUniFin
	}
	if wxname == "ie-fin" {
		return utils.URLWxIeFin
	}
	if wxname == "tan-money" {
		return utils.URLWxTanMoney
	}
	if wxname == "fin-circle" {
		return utils.URLWxFinCircle
	}
	return utils.URLWxFinCircle
}

func GrabWx(wxname string) []Essay {
	essays := []Essay{}
	url := matchUrl(wxname)
	if url == "" {
		return essays
	}

	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// handle error
		fmt.Println("err: ", err.Error())
	}
	htmlContent := string(body)
	hslice := strings.Split(htmlContent, "var msgList = ")
	jslice := strings.Split(hslice[1], "seajs.use")
	trimStr := strings.TrimSpace(jslice[0])
	realStr := trimStr[0 : len(trimStr)-1]

	result := gjson.GetMany(realStr,
		"list.#.app_msg_ext_info.title",       // 标题
		"list.#.app_msg_ext_info.content_url", // 地址
		"list.#.comm_msg_info.datetime",       // 时间
		"list.#.app_msg_ext_info.author")      // 作者
	if len(result) != 4 {
		return essays
	}
	itemNum := gjson.Get(result[0].String(), "#").Num
	for i := 0; i < int(itemNum); i++ {
		idx := fmt.Sprintf("%d", i)
		res00 := gjson.Get(result[0].String(), idx).String()
		res01 := gjson.Get(result[1].String(), idx).String()
		res02 := gjson.Get(result[2].String(), idx).Int()
		res03 := gjson.Get(result[3].String(), idx).String()
		essays = append(essays, Essay{
			Title:  res00,
			Url:    fmt.Sprintf("%s%s", "https://mp.weixin.qq.com", res01),
			Time:   time.Unix(res02, 0).Format("2006-01-02 15:04:05"),
			Author: res03,
		})
	}
	return essays
}
