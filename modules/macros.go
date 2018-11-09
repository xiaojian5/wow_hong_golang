package modules

import (
	sql "github.com/illidan33/sql-builder"
	"github.com/jmoiron/sqlx"
)

// 宏文本-数据库映射
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
	IsVerify int `json:"isVerify" db:"is_verify"`
}

// 快速创建
type SequenceMacro struct {
	// 名称
	SkillName string `json:"skillName"`
	// 优先级
	Level int `json:"level"`
	// 冷却时间(秒*100)
	Cooldown int `json:"cooldown"`
	// 剩余时间
	CDTime int `json:"-"`
}

// 获取模板列表
func GetMacroList(macro Macro) ([]Macro, error) {
	macroText := macro.Macro
	macro.Macro = "" // WhereByStruct跳过macro的条件

	builder := sql.Select("macros")
	builder.WhereByStruct(macro, true)
	if macroText != "" {
		builder.WhereLike("macro", macroText)
	}

	macros := make([]Macro, 0)

	conn := GetDbConn()
	rows, err := conn.Query(builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	sqlx.StructScan(rows, &macros)
	return macros, nil
}

func CreateMacro(macro Macro) error {
	builder := sql.Insert("macros")
	builder.InsertByStruct(macro)

	conn := GetDbConn()

	stmt, err := conn.Prepare(builder.String())
	if err != nil {
		return err
	}
	_, err = stmt.Exec(builder.Args()...)
	if err != nil {
		return err
	}

	return nil
}

func UpdateMacroByID(macro Macro, id int) error {
	builder := sql.Update("macros")
	builder.UpdateByStruct(macro, true)
	builder.WhereEq("id", id)

	conn := GetDbConn()

	_, err := conn.Exec(builder.String(), builder.Args()...)
	if err != nil {
		return err
	}

	return nil
}

func CreateSequence(temps []SequenceMacro) (macroText []string, maxTime int) {
	for _, value := range temps {
		if maxTime == 0 {
			maxTime = value.Cooldown
		} else {
			if value.Cooldown > maxTime {
				maxTime = value.Cooldown
			}
		}
	}

	for i := 0; i < maxTime; i++ {
		coolIndex := 0
		coolLevel := 0
		for j, value := range temps {
			if value.CDTime != 0 {
				continue
			}
			if coolLevel == 0 {
				coolLevel = value.Level
				coolIndex = j
			} else {
				if coolLevel > value.Level {
					coolLevel = value.Level
					coolIndex = j
				}
			}
		}
		if coolLevel != 0 {
			macroText = append(macroText, temps[coolIndex].SkillName)
			temps[coolIndex].CDTime = temps[coolIndex].Cooldown
		}
		// CD时间减1
		for k, value := range temps {
			if value.CDTime != 0 {
				temps[k].CDTime -= 1
			}
		}
	}

	return macroText, maxTime
}
