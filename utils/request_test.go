package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/tidwall/gjson"
)

func TestHuXiu(t *testing.T) {
	// curl 'https://wwww.huxiu.com/channel/ajaxGetMore'
	// --data 'huxiu_hash_code=217a7d7fac22ae0af7c7c8cb11968556&page=2&catId=102&last_dateline='
	// Content-Type: application/x-www-form-urlencoded

	const huxiuApi = "https://wwww.huxiu.com/channel/ajaxGetMore"
	client := &http.Client{}

	req, err := http.NewRequest("POST", huxiuApi, strings.NewReader("huxiu_hash_code=217a7d7fac22ae0af7c7c8cb11968556&page=2&catId=102&last_dateline="))
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
	value := gjson.Get(string(body), "data.data")
	fmt.Println("value =>", value)
}
