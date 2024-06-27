package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("Hello World")
	//创建一个gin应用
	r := gin.Default()
	//创建路由
	r.GET("/index", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
	r.Run(":8080")

}
