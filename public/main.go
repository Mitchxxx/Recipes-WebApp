package main

import "github.com/gin-gonic/gin"



func IndexHandler(c *gin.Context) {
	c.File("index.html")
}

func main () {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.GET("/", IndexHandler)
	router.Run()
}