package sql_builder

import (
	"fmt"
	"log"
	"reflect"
)

type SqlType uint8
type WhereType uint8

const (
	SQL_TYPE        SqlType = iota // unknown
	SQL_TYPE_INSERT                // insert
	SQL_TYPE_UPDATE                // update
	SQL_TYPE_SELECT                // select
	SQL_TYPE_DELETE                // delete
)

const (
	WHERE_TYPE         WhereType = iota // unknown
	WHERE_TYPE_EQ                       // equal
	WHERE_TYPE_NEQ                      // not equal
	WHERE_TYPE_GT                       // greater
	WHERE_TYPE_GTE                      // greater and equal
	WHERE_TYPE_LT                       // less
	WHERE_TYPE_LTE                      // less and equal
	WHERE_TYPE_Like                     // % before and after Like
	WHERE_TYPE_LikeBEF                  // % before Like
	WHERE_TYPE_LikeAFT                  // % after Like
)

type SqlBuilder struct {
	// table name
	table string
	// sql type
	sqlType SqlType
	// insert sql string
	insertStr string
	// update sql string
	updateStr string
	// where conditions
	whereStr string
	// args for db.Exec
	args []interface{}
}

// struct for build sql string of 'or'
type WhereOrCondition struct {
	// 字段对应数据库名称
	FieldName string
	// where条件
	WhereType WhereType
	// 字段值
	FieldValue interface{}
}

// init
func (build *SqlBuilder) Init(tableName string, sqlType SqlType) {
	build.insertStr = ""
	build.updateStr = ""
	build.whereStr = ""
	build.table = tableName
	build.sqlType = sqlType
}

// Build sql string of 'where' with one field
func (build *SqlBuilder) Where(fieldName string, whereType WhereType, fieldValue interface{}) {
	var conditionStr string
	conditionStr, fieldValue = build.buildWhereCondition(fieldName, whereType, fieldValue)

	if build.whereStr == "" {
		build.whereStr = conditionStr
	} else {
		build.whereStr = fmt.Sprintf("%s AND %s", build.whereStr, conditionStr)
	}
	build.args = append(build.args, fieldValue)
}

// Build sql string of 'in' condition with conditions
func (build *SqlBuilder) WhereIn(fieldName string, args []interface{}) {
	length := len(args)
	if length < 1 {
		log.Fatalf("Need args")
	}

	condition := ""
	for i := 0; i < length; i++ {
		if condition == "" {
			condition = "?"
			continue
		}
		condition = fmt.Sprintf("%s,?", condition)
	}

	if build.whereStr == "" {
		build.whereStr = fmt.Sprintf("%s IN (%s)", fieldName, condition)
	} else {
		build.whereStr = fmt.Sprintf("%s AND %s IN (%s)", build.whereStr, fieldName, condition)
	}

	for _, value := range args {
		build.args = append(build.args, value)
	}
}

// Build sql string of 'or' condition with struct WhereOrCondition
func (build *SqlBuilder) WhereOr(args []WhereOrCondition) {
	var orStr string
	for _, value := range args {
		conditionStr, fieldValue := build.buildWhereCondition(value.FieldName, value.WhereType, value.FieldValue)
		if orStr == "" {
			orStr = conditionStr
		} else {
			orStr = fmt.Sprintf("%s AND %s", orStr, conditionStr)
		}
		build.args = append(build.args, fieldValue)
	}
	if build.whereStr == "" {
		build.whereStr = orStr
	} else {
		build.whereStr = fmt.Sprintf("%s OR (%s)", build.whereStr, orStr)
	}
}

// Build sql string of update with one condition
func (build *SqlBuilder) UpdateSet(fieldName string, fieldValue interface{}) {
	if build.sqlType != SQL_TYPE_UPDATE {
		log.Fatalf("type error")
	}

	if build.updateStr == "" {
		build.updateStr = fmt.Sprintf("%s=?", fieldName)
	} else {
		build.updateStr = fmt.Sprintf("%s,%s=?", build.updateStr, fieldName)
	}
	build.args = append(build.args, fieldValue)
}

// Build sql string of update by struct
// The value of struct must has tag "db",what map to field of database.
func (build *SqlBuilder) UpdateByStruct(tableMap interface{}, skipEmpty bool) {
	if build.sqlType != SQL_TYPE_UPDATE {
		log.Fatalf("SQL type error")
	}

	tableType := reflect.TypeOf(tableMap)
	tableValue := reflect.ValueOf(tableMap)

	num := tableType.NumField()

	var sqlStr string
	for i := 0; i < num; i++ {
		dbTag := tableType.Field(i).Tag.Get("db")
		if dbTag == "" {
			log.Fatalf("Struct need tag 'db'")
		}
		value := tableValue.Field(i).Interface()
		if skipEmpty == true && (value == 0 || value == "") {
			continue
		}
		if sqlStr == "" {
			sqlStr = fmt.Sprintf("%s=?", dbTag)
		} else {
			sqlStr = fmt.Sprintf("%s,%s=?", sqlStr, dbTag)
		}
		build.args = append(build.args, value)
	}
	build.insertStr = fmt.Sprintf("UPDATE %s SET %s", build.table, sqlStr)
}

// Build sql string of insert by struct
// The value of struct must has tag "db",what map to field of database.
func (build *SqlBuilder) InsertByStruct(tableMap interface{}) {
	if build.sqlType != SQL_TYPE_INSERT {
		log.Fatalf("SQL type error")
	}

	tableType := reflect.TypeOf(tableMap)
	tableValue := reflect.ValueOf(tableMap)

	num := tableType.NumField()

	var sqlStr string
	var valStr string
	for i := 0; i < num; i++ {
		dbTag := tableType.Field(i).Tag.Get("db")
		if dbTag == "" {
			log.Fatalf("Struct need tag 'db'")
		}
		if sqlStr == "" {
			sqlStr = fmt.Sprintf("%s", dbTag)
			valStr = fmt.Sprintf("?")
		} else {
			sqlStr = fmt.Sprintf("%s,%s", sqlStr, dbTag)
			valStr = fmt.Sprintf("%s,?", valStr)
		}
		build.args = append(build.args, tableValue.Field(i).Interface())
	}
	build.insertStr = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", build.table, sqlStr, valStr)
}

// Get sql string
func (build *SqlBuilder) String() string {
	whereStr := ""
	if build.whereStr != "" {
		whereStr = fmt.Sprintf("WHERE %s", build.whereStr)
	}
	switch build.sqlType {
	case SQL_TYPE_INSERT:
		return build.insertStr
	case SQL_TYPE_UPDATE:
		return fmt.Sprintf("UPDATE %s SET %s %s", build.table, build.updateStr, whereStr)
	case SQL_TYPE_SELECT:
		return fmt.Sprintf("SELECT * FROM %s %s", build.table, whereStr)
	case SQL_TYPE_DELETE:
		return fmt.Sprintf("DELETE FROM %s %s", build.table, whereStr)
	case SQL_TYPE:
		log.Fatalf("sql type error")
	}
	return ""
}

// Get all Args
func (build *SqlBuilder) Args() []interface{} {
	return build.args
}

// Build condition string
func (build *SqlBuilder) buildWhereCondition(fieldName string, whereType WhereType, fieldValue interface{}) (string, interface{}) {
	if whereType == WHERE_TYPE {
		log.Fatalf("Where type error")
	}

	condition := ""
	switch whereType {
	case WHERE_TYPE_EQ:
		condition = "="
		break
	case WHERE_TYPE_NEQ:
		condition = "<>"
		break
	case WHERE_TYPE_GT:
		condition = ">"
		break
	case WHERE_TYPE_GTE:
		condition = ">="
		break
	case WHERE_TYPE_LT:
		condition = "<"
		break
	case WHERE_TYPE_LTE:
		condition = "<="
		break
	case WHERE_TYPE_Like:
		fieldValue = fmt.Sprintf("%%%s%%", fieldValue)
		condition = " LIKE "
		break
	case WHERE_TYPE_LikeAFT:
		fieldValue = fmt.Sprintf("%s%%", fieldValue)
		condition = " LIKE "
		break
	case WHERE_TYPE_LikeBEF:
		fieldValue = fmt.Sprintf("%%%s", fieldValue)
		condition = " LIKE "
		break
	default:
		break
	}

	return fmt.Sprintf("%s%s?", fieldName, condition), fieldValue
}
