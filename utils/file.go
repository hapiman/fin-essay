package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func ReadAllOnce(path string) string {
	b, e := ioutil.ReadFile("d:/goTest/123.txt")
	if e != nil {
		fmt.Println("read file error")
		return ""
	}
	return string(b)
}

func ReadLineEachTime(path string) []string {
	cs := []string{}
	fi, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return cs
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		cs = append(cs, string(a))
	}
	return cs
}
