package modules

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var DbConn *sqlx.DB

func init() {
	var err error
	DbConn, err = sqlx.Open("mysql", "test:test@tcp(127.0.0.1:3306)/wow_hong?charset=utf8")
	CheckErr("Connect Database", err)

	DbConn.SetMaxOpenConns(2000)
	DbConn.SetMaxIdleConns(1000)
	err = DbConn.Ping()
	CheckErr("Ping Database", err)
}

func GetDbConn() *sqlx.DB {
	var err error
	err = DbConn.Ping()
	if err != nil {
		DbConn, err = sqlx.Open("mysql", "test:test@tcp(127.0.0.1:3306)/wow_hong?charset=utf8")
		CheckErr("Connect Database", err)

		DbConn.SetMaxOpenConns(2000)
		DbConn.SetMaxIdleConns(1000)

		err = DbConn.Ping()
		if err == nil {
			return DbConn
		}
		CheckErr("Ping Database", err)
	}
	return DbConn
}

func CheckErr(msg string, err error) {
	if err != nil {
		fmt.Fprintf(gin.DefaultWriter, "%s : %+v\n", msg, err)
	}
}
