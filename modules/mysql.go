package modules

import (
	"github.com/sirupsen/logrus"
	"os"
	"github.com/jmoiron/sqlx"
	"log"
	"fmt"
	"reflect"
)

type SqlType uint8

const (
	SQL_TYPE        SqlType = iota
	SQL_TYPE_INSERT
	SQL_TYPE_UPDATE
	SQL_TYPE_SELECT
	SQL_TYPE_DELETE
)

var DbConn *sqlx.DB

type BuildSql struct {
	// 表名
	Table string `json:"table"`
	// sql操作类型
	Type SqlType `json:"sqlType"`
	// 创建sql
	InsertStr string `json:"insertStr"`
	// 更新sql
	UpdateStr string `json:"sqlStr"`
	// 操作条件
	WhereStr string `json:"whereStr"`
	// 操作数据
	Args []interface{} `json:"args"`
}

func init() {
	var err error
	DbConn, err = sqlx.Open("mysql", "root:xiaohu@tcp(127.0.0.1:3306)/wow_hong?charset=utf8")
	if err != nil {
		logrus.Errorf("connect error: %s", err.Error())
		os.Exit(0)
	}
	DbConn.SetMaxOpenConns(2000)
	DbConn.SetMaxIdleConns(1000)
	DbConn.Ping()
}

func CheckErr(msg string, err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化
func (build *BuildSql) Init(tableName string, sqlType SqlType) BuildSql {
	build.Table = tableName
	build.Type = sqlType
	return *build
}

// 构建where语句
func (build *BuildSql) Eq(fieldName string, fieldValue interface{}) {
	if fieldValue == "" || fieldValue == 0 {

	} else {
		if build.WhereStr == "" {
			build.WhereStr = fmt.Sprintf("WHERE %s=?", fieldName)
		} else {
			build.WhereStr = fmt.Sprintf("%s AND %s=?", build.WhereStr, fieldName)
		}
		build.Args = append(build.Args, fieldValue)
	}
}

// 构建where语句
func (build *BuildSql) Like(fieldName string, fieldValue interface{}) {
	if fieldValue == "" || fieldValue == 0 {

	} else {
		if build.WhereStr == "" {
			build.WhereStr = fmt.Sprintf("WHERE %s LIKE ?", fieldName)
		} else {
			build.WhereStr = fmt.Sprintf("%s AND %s LIKE ?", build.WhereStr, fieldName)
		}
		build.Args = append(build.Args, fieldValue)
	}
}

// 构建update语句
func (build *BuildSql) Set(fieldName string, fieldValue interface{}) {
	if build.Type != SQL_TYPE_UPDATE {
		log.Fatalf("type error")
	}

	if build.UpdateStr == "" {
		build.UpdateStr = fmt.Sprintf("%s=?", fieldName)
	} else {
		build.UpdateStr = fmt.Sprintf("%s,%s=?", build.UpdateStr, fieldName)
	}
	build.Args = append(build.Args, fieldValue)
}

// 根据struct构建update语句
func (build *BuildSql) Update(tableMap interface{}) {
	if build.Type != SQL_TYPE_UPDATE {
		log.Fatalf("type error")
	}

	tableType := reflect.TypeOf(tableMap)
	tableValue := reflect.ValueOf(tableMap)

	num := tableType.NumField()

	var sqlStr string
	for i := 0; i < num; i++ {
		value := tableValue.Field(i).Interface()
		if value == 0 || value == "" {
			continue
		}
		if sqlStr == "" {
			sqlStr = fmt.Sprintf("%s=?", tableType.Field(i).Tag.Get("db"))
		} else {
			sqlStr = fmt.Sprintf("%s,%s=?", sqlStr, tableType.Field(i).Tag.Get("db"))
		}
		build.Args = append(build.Args, value)
	}
	build.InsertStr = fmt.Sprintf("UPDATE %s SET %s", build.Table, sqlStr)
}

// 构建插入语句
func (build *BuildSql) Insert(tableMap interface{}) {
	if build.Type != SQL_TYPE_INSERT {
		log.Fatalf("type error")
	}

	tableType := reflect.TypeOf(tableMap)
	tableValue := reflect.ValueOf(tableMap)

	num := tableType.NumField()

	var sqlStr string
	var valStr string
	for i := 0; i < num; i++ {
		if tableType.Field(i).Tag.Get("db") == "id" {
			continue
		}
		if sqlStr == "" {
			sqlStr = fmt.Sprintf("%s", tableType.Field(i).Tag.Get("db"))
			valStr = fmt.Sprintf("?")
		} else {
			sqlStr = fmt.Sprintf("%s,%s", sqlStr, tableType.Field(i).Tag.Get("db"))
			valStr = fmt.Sprintf("%s,?", valStr)
		}
		build.Args = append(build.Args, tableValue.Field(i).Interface())
	}
	build.InsertStr = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", build.Table, sqlStr, valStr)
}

func (build *BuildSql) String() string {
	switch build.Type {
	case SQL_TYPE_INSERT:
		return build.InsertStr
	case SQL_TYPE_UPDATE:
		return fmt.Sprintf("UPDATE %s SET %s %s", build.Table, build.UpdateStr, build.WhereStr)
	case SQL_TYPE_SELECT:
		return fmt.Sprintf("SELECT * FROM %s %s", build.Table, build.WhereStr)
	case SQL_TYPE_DELETE:
		return fmt.Sprintf("DELETE FROM %s %s", build.Table, build.WhereStr)
	case SQL_TYPE:
		log.Fatalf("sql type error")
	}
	return ""
}
