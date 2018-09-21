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
	"strings"
	"math"
)

var (
	port int
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			errNew := err.(error)
			modules.CheckErr("Error", errNew)
		}
	}()

	flag.IntVar(&port, "port", 8000, "listen port")
	flag.Parse()

	//gin.SetMode(gin.ReleaseMode)
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
	router.PUT("/macros/:id", UpdateMacro)
	router.POST("/log/:method", CreateLoginLog)
	router.POST("/create-sequence", CreateSequence)
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/html/index.html")
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
	isVerify, err := strconv.Atoi(c.DefaultQuery("isVerify", "1"))
	modules.CheckErr("isVerify", err)

	macro := modules.Macro{
		ID:           id,
		ProfessionID: professionId,
		MasteryID:    masteryId,
		Macro:        c.DefaultQuery("macro", ""),
		IsVerify:     isVerify,
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
	} else {
		macro.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
		macro.IsVerify = 2

		result := modules.CreateMacro(macro)
		if result == true {
			c.JSON(http.StatusOK, gin.H{})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{})
		}
	}
}

func UpdateMacro(c *gin.Context) {
	token, err := c.Cookie("token")
	// 简单的权限验证
	if err != nil || token != "test" {
		c.JSON(http.StatusForbidden, gin.H{})
	} else {
		id, err := strconv.Atoi(c.Param("id"))
		modules.CheckErr("id", err)

		macro := modules.Macro{
			ID:         id,
			UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		}

		err = c.BindJSON(&macro)
		modules.CheckErr("BindJSON", err)

		result := modules.UpdateMacroByID(macro, macro.ID)
		c.JSON(http.StatusOK, gin.H{"result": result,})
	}
}

func CreateLoginLog(c *gin.Context) {
	method := c.Param("method")
	// 记录日志
	ip := c.ClientIP()
	go modules.CreateLog(ip, method)

	c.JSON(http.StatusOK, gin.H{})
}

func CreateSequence(c *gin.Context) {
	temps := make([]modules.SequenceMacro, 2)

	body, err := ioutil.ReadAll(c.Request.Body)
	modules.CheckErr("CreateSequence", err)
	err = json.Unmarshal(body, &temps)
	modules.CheckErr("CreateSequence unmarshal", err)

	for i, value := range temps {
		if value.Cooldown == 0 {
			temps[i].Cooldown = 100
		}
		temps[i].SkillName = strings.Replace(value.SkillName, " ", "", -1)
		temps[i].SkillName = strings.Replace(value.SkillName, "\n", "", -1)
	}
	fmt.Fprintf(gin.DefaultWriter, "%s : %+v\n", "test", temps)

	macros, maxTime := modules.CreateSequence(temps)
	maxTime = int(math.Ceil(float64(maxTime) / 100))
	macroText := fmt.Sprintf("#showtooltip <br>/castsequence reset=%d %s", maxTime, strings.Join(macros, ","));
	c.JSON(http.StatusOK, gin.H{"text": macroText, "desc": fmt.Sprintf("- 最后一次按键【%d】秒后，将重置<br>- 技能按照左侧循环，时间和技能顺序可以自己修改！", maxTime)})
}
