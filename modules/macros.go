package modules

import (
	"github.com/jmoiron/sqlx"
	"fmt"
)

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
	// 专精
	MasteryID int `json:"masteryId" db:"mastery_id"`
	// 职业
	ProfessionID int `json:"professionId" db:"profession_id"`
	// 是否通过审核
	IsVerify int8 `json:"isVerify" db:"is_verify"`
}

// 获取模板列表
func GetMacroList(macro Macro) []Macro {
	var build BuildSql
	build.Init("macros", SQL_TYPE_SELECT)

	build.Eq("mastery_id", macro.MasteryID)
	build.Eq("profession_id", macro.ProfessionID)
	build.Eq("id", macro.ID)
	build.Like("macro", fmt.Sprintf("%%%s%%", macro.Macro))

	macros := make([]Macro, 0)

	rows, err := DbConn.Query(build.String(), build.Args...)
	CheckErr("GetMacroList", err)

	sqlx.StructScan(rows, &macros)
	return macros
}

func CreateMacro(macro Macro) bool {
	var build BuildSql
	build.Init("macros", SQL_TYPE_INSERT)
	build.Insert(macro)

	tx, err := DbConn.Begin()
	if err != nil {
		CheckErr("CreateMacro beginx", err)
		return false
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(build.String())
	_, err = stmt.Exec(build.Args...)
	if err != nil {
		tx.Rollback()
		CheckErr("CreateMacro exec", err)
	}

	err = tx.Commit()
	if err != nil {
		CheckErr("CreateMacro commit", err)
		return false
	}

	return true
}

func UpdateMacro(macro Macro) bool {
	var build BuildSql
	build.Init("macros", SQL_TYPE_UPDATE)
	build.Update(macro)

	tx, err := DbConn.Beginx()
	if err != nil {
		CheckErr("CreateMacro beginx", err)
		return false
	}

	_, err = tx.Exec(build.String(), build.Args...)
	if err != nil {
		tx.Rollback()
		CheckErr("CreateMacro exec", err)
		return false
	}

	err = tx.Commit()
	if err != nil {
		CheckErr("CreateMacro commit", err)
		return false
	}

	return true
}
