package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hapiman/fin-essay/crawer"
)

func formatSuccResp(iList []crawer.Essay) gin.H {
	return gin.H{
		"data": map[string]interface{}{
			"list": iList,
		},
		"success": true,
	}
}

func main() {
	r := gin.Default()
	go crawer.StartTaskRobot()
	r.GET("/fin/iyiou", func(c *gin.Context) {
		iList := crawer.ReadEssay("iyiou")
		c.JSON(200, formatSuccResp(iList))
	})
	r.GET("/fin/wdzj", func(c *gin.Context) {
		iList := crawer.ReadEssay("wdzj")
		c.JSON(200, formatSuccResp(iList))
	})
	r.GET("/fin/huxiu", func(c *gin.Context) {
		iList := crawer.Grab_HuXiu()
		c.JSON(200, formatSuccResp(iList))
	})
	r.GET("/fin/wx", func(c *gin.Context) {
		iList := crawer.GrabWx("")
		c.JSON(200, formatSuccResp(iList))
	})
	r.GET("/fin/wx/:wxname", func(c *gin.Context) {
		wxname := c.Param("wxname")
		iList := crawer.GrabWx(wxname)
		c.JSON(200, formatSuccResp(iList))
	})
	r.Run()
}
