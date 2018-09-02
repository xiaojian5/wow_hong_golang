package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

// 职业
type Profession struct {
	ID int `json:"id"`
	// 父级ID
	PID int `json:"pid"`
	// 名称
	Name string `json:"name"`
}

func GetProfessionList(c *gin.Context) {
	db, err := sqlx.Connect("mysql", "hong:hong486@tcp(45.32.254.107:3306)/wow_hong")
	CheckError(c, err)
	defer db.Close()

	pid := c.Query("pid")

	if pid == "" {
		pid = "0"
	}

	macros := make([]Profession, 0)

	err = db.Select(&macros, "SELECT id,pid,name from profession where pid=?", pid)
	CheckError(c, err)

	c.JSON(http.StatusOK, macros)
}