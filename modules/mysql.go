package modules

import (
	"os"
	"github.com/jmoiron/sqlx"
	"log"
	"fmt"
	"github.com/gin-gonic/gin"
)

var DbConn *sqlx.DB

func init() {
	var err error
	DbConn, err = sqlx.Open("mysql", "test:test@tcp(127.0.0.1:3306)/wow_hong?charset=utf8")
	if err != nil {
		log.Fatalf("connect error: %s", err.Error())
		os.Exit(0)
	}
	DbConn.SetMaxOpenConns(2000)
	DbConn.SetMaxIdleConns(1000)
	DbConn.Ping()
}

func CheckErr(msg string, err error) {
	if err != nil {
		fmt.Fprintf(gin.DefaultWriter, "%s : %+v\n", msg, err)
		log.Fatal(err)
	}
}
