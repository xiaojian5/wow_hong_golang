package modules

import (
	"github.com/jmoiron/sqlx"
	sql "github.com/illidan33/sql-builder"
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
	builder := sql.Select("macros")

	if macro.MasteryID != 0 {
		builder.WhereEq("mastery_id", macro.MasteryID)
	}
	if macro.ProfessionID != 0 {
		builder.WhereEq("profession_id", macro.ProfessionID)
	}
	if macro.ID != 0 {
		builder.WhereEq("id", macro.ID)
	}
	if macro.Macro != "" {
		builder.WhereLike("macro", macro.Macro)
	}

	macros := make([]Macro, 0)

	rows, err := DbConn.Query(builder.String(), builder.Args()...)
	CheckErr("GetMacroList", err)

	sqlx.StructScan(rows, &macros)
	return macros
}

func CreateMacro(macro Macro) bool {
	builder := sql.Insert("macros")
	builder.InsertByStruct(macro)

	tx, err := DbConn.Begin()
	if err != nil {
		CheckErr("CreateMacro beginx", err)
		return false
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(builder.String())
	_, err = stmt.Exec(builder.Args()...)
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
	builder := sql.Update("macros")
	builder.UpdateByStruct(macro, true)

	tx, err := DbConn.Beginx()
	if err != nil {
		CheckErr("CreateMacro beginx", err)
		return false
	}

	_, err = tx.Exec(builder.String(), builder.Args()...)
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
