package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"gotest/wow_hong_v2_service/modules"
	"github.com/jmoiron/sqlx"
	"log"
)

func main() {
	/*dbConfig := modules.Config{
		Host: "45.32.254.107",
		Port: "3306",
		UserName: "hong",
		Pwd: "hong486",
	}*/
	db, err := sqlx.Open("mysql", "hong:hong486@tcp(45.32.254.107:3306)/wow_hong?parseTime=true")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	router := gin.Default()

	router.GET("/", IndexFunc)
	router.GET("/macroList", modules.GetMacroList)
	router.GET("/professionList", modules.GetProfessionList)

	router.Run(":8080")
}

func IndexFunc(c *gin.Context) {
	c.String(http.StatusOK, "hello world!")
}
