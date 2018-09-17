# sql-builder
A tool to create sql string for golang

# 安装

```
go get -u github.com/illidan33/sql-builder
```

# 实例

### 查询

1. 预先定义包含db标签的struct，自动设置查询字段为包含在struct中的db标签

```
type Skill struct {
	Condition string `json:"condition" db:"condition"`
	Desc      string `json:"desc" db:"desc"`
	SkillType int    `json:"skillType" db:"skill_type"`
}

skill := Skill{
    Condition: "",
    SkillType: 1,
}

builder := Select("skill")
// SelectByStruct会设置查询字段，同时会设置查询条件
builder.SelectByStruct(skill, true) // 第二个参数为是否跳过空值

Dbconn.Query(builder.String(), builder.Args()...) // 放入数据库查询中

// sql：SELECT `condition`,`desc`,`skill_type` FROM `skill` WHERE `skill_type`=?;
```

```
type Skill struct {
	Condition string `json:"condition" db:"condition"`
	Desc      string `json:"desc" db:"desc"`
	SkillType int    `json:"skillType" db:"skill_type"`
}

skill := Skill{
    Condition: "",
    SkillType: 1,
}

builder := Select("skill")
// WhereByStruct只包含查询条件，不包含设置查询字段
builder.WhereByStruct(skill, true) // 第二个参数为是否跳过空值

Dbconn.Query(builder.String(), builder.Args()...) // 放入数据库查询中

// sql：SELECT * FROM `skill` WHERE `skill_type`=?;
```

2. 自定义查询条件

```
builder := Select("skill")
builder.SetSearchFields("`condition`,`skill_type`") // 不设置，默认为*
builder.WhereEq("skill_type", 1)
builder.WhereIn("skill_type", []interface{}{1, 2})
builder.WhereGt("skill_type", 1)
builder.WhereLt("skill_type", 1)
builder.WhereLike("condition", "vic")
builder.WhereOr([]WhereOrCondition{
    {
        FieldName:  "skill_type",
        WhereType:  WHERE_TYPE_EQ,
        FieldValue: 1,
    },
})
builder.GroupBy("`id`,`skill_type`")
builder.OrderBy("`skill_type` asc")
builder.Limit(0, 20)

Dbconn.Query(builder.String(), builder.Args()...) // 放入数据库查询中

// sql：SELECT `condition`,`skill_type` FROM `skill` WHERE (`skill_type`=? AND skill_type IN (?,?) AND `skill_type`>? AND `skill_type`<? AND `condition` LIKE ?) OR (`skill_type`=?) GROUP BY `id`,`skill_type` ORDER BY `skill_type` asc LIMIT 0,20;
// 需要自己注意group by/order by/limit的顺序
```

### 插入

```
skill := Skill{
    Condition: "test",
    Desc: "",
    SkillType: 1,
}

builder := Insert("skill")
builder.InsertByStruct(skill)

Dbconn.Query(builder.String(), builder.Args()...) // 放入数据库查询中

// sql：INSERT INTO `skill`(condition,desc,skill_type) VALUES(?,?,?);
```

### 修改

```
skill := Skill{
    Condition: "",
    Desc: "",
    SkillType: 1,
}

builder := Update("skill")
builder.UpdateByStruct(skill, true) // 第二个参数：是否跳过空值，如果为true，会跳过空字符串和0值
// 更新条件
builder.WhereEq("skill_type", 1)
builder.WhereIn("skill_type", []interface{}{1, 2})
builder.WhereGt("skill_type", 1)
builder.WhereLt("skill_type", 1)
builder.WhereLike("condition", "vic")
builder.WhereOr([]WhereOrCondition{
    {
        FieldName:  "skill_type",
        WhereType:  WHERE_TYPE_EQ,
        FieldValue: 1,
    },
})

Dbconn.Query(builder.String(), builder.Args()...) // 放入数据库查询中

// sql：UPDATE `skill` SET `skill_type`=? WHERE (`skill_type`=? AND skill_type IN (?,?) AND `skill_type`>? AND `skill_type`<? AND `condition` LIKE ?) OR (`skill_type`<? AND `skill_type`>?);
```

### 删除

```
skill := Skill{
    SkillType: 1,
}

builder := Delete("skill")
builder.WhereEq("skill_type", 1)
builder.WhereIn("skill_type", []interface{}{1, 2})
builder.WhereOr([]WhereOrCondition{
    {
        FieldName:  "skill_type",
        WhereType:  WHERE_TYPE_EQ,
        FieldValue: 1,
    },
})
builder.WhereGt("skill_type", 1)
builder.WhereLt("skill_type", 1)
builder.WhereLike("condition", "vic")

Dbconn.Query(builder.String(), builder.Args()...) // 放入数据库查询中
```
