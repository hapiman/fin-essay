package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestMountSubrouterOn(t *testing.T) {
	createdTime := "2019-07-20 16:01:23"
	create_Time, _ := time.ParseInLocation("2006-01-02 15:04:05", createdTime, time.Now().Location())
	sub := time.Now().Sub(create_Time)
	fmt.Println(sub.Hours())
}
