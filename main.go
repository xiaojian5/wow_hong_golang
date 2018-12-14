package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/illidan33/wow_hong_golang/modules"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	port     int
	tokenStr = "test"
)

func main() {
	envPath := os.Getenv("GOPATH")
	if envPath == "" {
		envPath = "/test";
	}
	rootPath := fmt.Sprintf("%s/src/github.com/illidan33/wow_hong_golang/", envPath);
	flag.IntVar(&port, "port", 8000, "listen port")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.LoadHTMLGlob(rootPath + "html/*")

	// 设置静态资源
	router.Static("/js", rootPath+"js")
	router.Static("/css", rootPath+"css")
	router.Static("/img", rootPath+"img")
	router.StaticFile("/favicon.ico", rootPath+"favicon.ico")
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status": 404,
			"error":  "404, page not exists!",
		})
	})

	// html
	router.GET("/", Index)
	router.GET("/html/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})
	router.GET("/Macro/:type", MacroIndex)

	// data
	router.GET("/macros", getMacroList)
	router.POST("/macros", CreateMacro)
	router.PUT("/macros", UpdateMacro)
	router.POST("/create-sequence", CreateSequence)

	router.Run(fmt.Sprintf(":%d", port))
}

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"page": "index",
	})
}
func MacroIndex(c *gin.Context) {
	c.Request.ParseForm()

	page := c.Param("type")

	if page != "" {
		ip := c.ClientIP()
		go CheckLoginLog(ip, page)
	}

	switch page {
	case "macroByHand":
		c.HTML(http.StatusOK, "macroByHand.html", gin.H{
			"page": page,
		})
		break
	case "macroShare":
		c.HTML(http.StatusOK, "macroShare.html", gin.H{
			"page": page,
		})
		break
	case "macroCtSequence":
		c.HTML(http.StatusOK, "macroCtSequence.html", gin.H{
			"page": page,
		})
		break
	case "macroInfo":
		c.HTML(http.StatusOK, "macroInfo.html", gin.H{
			"page": page,
		})
		break
	case "macroEditList":
		c.HTML(http.StatusOK, "macroEditList.html", gin.H{
			"page": page,
		})
		break
	case "macroList":
		c.HTML(http.StatusOK, "macroList.html", gin.H{
			"page": page,
		})
		break
	case "macroVerify":
		token, err := c.Cookie("token")
		// 简单的权限验证
		if err != nil || token != tokenStr {
			c.JSON(http.StatusForbidden, gin.H{
				"status": 403,
				"error":  "403, Forbidden!",
			})
		} else {
			c.HTML(http.StatusOK, "macroVerify.html", gin.H{
				"page": page,
			})
		}

		break
	case "macroCreate":
		c.HTML(http.StatusOK, "macroCreate.html", gin.H{
			"page": page,
		})
		break
	default:
		c.HTML(http.StatusOK, "index.html", gin.H{
			"page": "index",
		})
	}
}

func getMacroList(c *gin.Context) {
	c.Request.ParseForm()

	id, _ := strconv.Atoi(c.DefaultQuery("id", ""))
	professionId, _ := strconv.Atoi(c.DefaultQuery("professionId", ""))
	masteryId, _ := strconv.Atoi(c.DefaultQuery("masteryId", ""))
	isVerify, _ := strconv.Atoi(c.DefaultQuery("isVerify", ""))

	macro := modules.Macro{
		ID:           id,
		ProfessionID: professionId,
		MasteryID:    masteryId,
		Macro:        c.DefaultQuery("macro", ""),
		IsVerify:     isVerify,
	}
	result, err := modules.GetMacroList(macro)
	if err != nil {
		modules.CheckErr("getMacroList", err)
		c.JSON(http.StatusBadRequest, gin.H{})
	} else {
		c.JSON(http.StatusOK, result)
	}
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

		err := modules.CreateMacro(macro)
		if err != nil {
			modules.CheckErr("createMacro", err)
			c.JSON(http.StatusBadRequest, gin.H{})
		} else {
			c.JSON(http.StatusOK, gin.H{})
		}
	}
}

func UpdateMacro(c *gin.Context) {
	token, err := c.Cookie("token")
	// 简单的权限验证
	if err != nil || token != tokenStr {
		c.JSON(http.StatusForbidden, gin.H{})
	} else {
		macro := modules.Macro{
			UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
		}

		err = c.BindJSON(&macro)
		if err != nil {
			modules.CheckErr("BindJSON", err)
		}

		err := modules.UpdateMacroByID(macro, macro.ID)
		if err != nil {
			modules.CheckErr("updateMacro", err)
			c.JSON(http.StatusOK, gin.H{"result": false,})
		} else {
			c.JSON(http.StatusOK, gin.H{"result": true,})
		}
	}
}

func CreateSequence(c *gin.Context) {
	c.Request.ParseForm()
	temps := make([]modules.SequenceMacro, 2)

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		modules.CheckErr("CreateSequence", err)
	}
	err = json.Unmarshal(body, &temps)
	if err != nil {
		modules.CheckErr("CreateSequence unmarshal", err)
	}

	for i, value := range temps {
		if value.Cooldown == 0 {
			temps[i].Cooldown = 100
		}
		temps[i].SkillName = strings.Replace(value.SkillName, " ", "", -1)
		temps[i].SkillName = strings.Replace(value.SkillName, "\n", "", -1)
	}

	macros, maxTime := modules.CreateSequence(temps)
	maxTime = int(math.Ceil(float64(maxTime) / 100))
	macroText := fmt.Sprintf("#showtooltip <br>/castsequence reset=%d %s", maxTime, strings.Join(macros, ","));
	c.JSON(http.StatusOK, gin.H{"text": macroText, "desc": fmt.Sprintf("- 最后一次按键【%d】秒后，将重置<br>- 技能按照左侧循环，时间和技能顺序可以自己修改！", maxTime)})
}

func CheckLoginLog(ip string, method string) {
	date := time.Now().Format("2006-01-02")

	log, err := modules.GetLog(ip, method, date)
	if err != nil || log.ID == 0 {
		err = modules.CreateLog(ip, method)
		if err != nil {
			modules.CheckErr("CreateLog", err)
		}
	} else {
		err = modules.UpdateLog(log.ID, log.Count+1)
		if err != nil {
			modules.CheckErr("UpdateLog", err)
		}
	}
}
