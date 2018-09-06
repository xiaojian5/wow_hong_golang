package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gotest/gin_test/modules"
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"time"
	"fmt"
)

func main() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// 设置静态资源
	router.Static("/js", "js")
	router.Static("/css", "css")
	router.Static("/html", "html")
	router.StaticFile("/favicon.ico", "./favicon.ico")

	router.GET("/macros", getMacroList)
	router.POST("/macros", CreateMacro)

	router.Run(":8000")
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
	}
	result := modules.GetMacroList(macro)

	c.JSON(http.StatusOK, result)
}

func CreateMacro(c *gin.Context) {
	macro := modules.Macro{}

	body, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Fprintf(gin.DefaultWriter, "%s\n", string(body))
	err := json.Unmarshal(body, &macro)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	macro.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(gin.DefaultWriter, "%+v\n", macro)
	result := modules.CreateMacro(macro)
	if result == true {
		c.JSON(http.StatusOK, gin.H{})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{})
	}
}
