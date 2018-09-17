package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"wow_hong_golang/modules"
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"time"
	"fmt"
	"flag"
	"log"
)

var (
	port int
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	flag.IntVar(&port, "port", 8000, "listen port")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 设置静态资源
	router.Static("/js", "js")
	router.Static("/css", "css")
	router.Static("/html", "html")
	router.StaticFile("/favicon.ico", "./favicon.ico")
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": 404,
			"error":  "404, page not exists!",
		})
	})

	router.GET("/macros", getMacroList)
	router.POST("/macros", CreateMacro)
	router.PUT("/macros", UpdateMacro)
	router.GET("/", func(c *gin.Context) {
		// 记录日志
		ip := c.ClientIP()
		modules.CreateLog(ip, "index")

		c.Redirect(http.StatusMovedPermanently, "/html/")
	})

	router.Run(fmt.Sprintf(":%d", port))
}

func getMacroList(c *gin.Context) {
	id, err := strconv.Atoi(c.DefaultQuery("id", "0"))
	modules.CheckErr("id", err)
	professionId, err := strconv.Atoi(c.DefaultQuery("professionId", "0"))
	modules.CheckErr("professionId", err)
	masteryId, err := strconv.Atoi(c.DefaultQuery("masteryId", "0"))
	modules.CheckErr("masteryId", err)

	macro := modules.Macro{
		ID:           id,
		ProfessionID: professionId,
		MasteryID:    masteryId,
		Macro:        c.DefaultQuery("macro", ""),
	}
	result := modules.GetMacroList(macro)

	c.JSON(http.StatusOK, result)
}

func CreateMacro(c *gin.Context) {
	macro := modules.Macro{}

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &macro)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	macro.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	result := modules.CreateMacro(macro)
	if result == true {
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
}

func UpdateMacro(c *gin.Context) {
	macro := modules.Macro{}

	body, _ := ioutil.ReadAll(c.Request.Body)
	err := json.Unmarshal(body, &macro)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	macro.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	result := modules.UpdateMacro(macro)
	if result == true {
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
}
