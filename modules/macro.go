package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
)

// 宏
type Macro struct {
	ID int `json:"id" db:"id"`
	// 标题
	Title string `json:"title" db:"title"`
	// 模板
	Macro string `json:"macro" db:"macro"`
	// 更新时间
	UpdateTime string `json:"updateTime" db:"updatetime"`
	// 作者
	Author string `json:"author" db:"author"`
	// 专精ID
	MasteryID int `json:"masteryId" db:"mastery_id"`
	// 职业ID
	ProfessionID int `json:"professionId" db:"profession_id"`
	// 是否通过审核
	IsVerify int8 `json:"isVerify" db:"is_verify"`
}

// 获取宏列表
func GetMacroList(c *gin.Context) {
	db, err := sqlx.Connect("mysql", "hong:hong486@tcp(45.32.254.107:3306)/wow_hong")
	CheckError(c, err)
	defer db.Close()

	materyId := c.Query("materyId")
	professionId := c.Query("professionId")
	isVerify := c.Query("isVerify")

	if materyId == "" {
		materyId = "0"
	}
	if professionId == "" {
		professionId = "1"
	}
	if isVerify == "" {
		isVerify = "1"
	}

	macros := make([]Macro, 0)

	err = db.Select(&macros, "SELECT id,title,macro,updatetime,author,mastery_id,profession_id,is_verify from macros where is_verify=? and mastery_id=? and profession_id=?", isVerify, materyId, professionId)
	CheckError(c, err)

	c.JSON(http.StatusOK, macros)
}
