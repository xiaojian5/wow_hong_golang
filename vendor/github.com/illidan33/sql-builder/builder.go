package sql_builder

import (
	"fmt"
	"log"
	"reflect"
)

type SqlType string
type WhereType uint8

const (
	SQL_TYPE_INSERT SqlType = "INSERT" // insert
	SQL_TYPE_UPDATE SqlType = "UPDATE" // update
	SQL_TYPE_SELECT SqlType = "SELECT" // select
	SQL_TYPE_DELETE SqlType = "DELETE" // delete
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

type sqlBuilder struct {
	// 表名/table name
	table string
	// 查询字段/select fields
	fields string
	// sql type
	sqlType SqlType
	// sql string
	handleStr string
	// where conditions
	whereStr string
	// 占位符/sql flag
	flag string
	// args for db.Exec
	args []interface{}
}

type whereBuilder struct {
	sqlBuilder
}

type SelectSqlBuilder struct {
	whereBuilder
}
type UpdateSqlBuilder struct {
	whereBuilder
}
type InsertSqlBuilder struct {
	sqlBuilder
}
type DeleteSqlBuilder struct {
	whereBuilder
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
func Select(tableName string) *SelectSqlBuilder {
	build := &SelectSqlBuilder{}
	build.init(tableName, SQL_TYPE_SELECT)

	return build
}

// init
func Update(tableName string) *UpdateSqlBuilder {
	build := &UpdateSqlBuilder{}
	build.init(tableName, SQL_TYPE_UPDATE)

	return build
}

// init
func Insert(tableName string) *InsertSqlBuilder {
	build := &InsertSqlBuilder{}
	build.init(tableName, SQL_TYPE_INSERT)

	return build
}

// init
func Delete(tableName string) *DeleteSqlBuilder {
	build := &DeleteSqlBuilder{}
	build.init(tableName, SQL_TYPE_DELETE)

	return build
}

// init
func (build *sqlBuilder) init(tableName string, sqlType SqlType) *sqlBuilder {
	build.sqlType = sqlType
	build.table = tableName
	build.flag = "?"
	build.fields = "*"

	return build
}

// Build sql string of update by struct
// The value of struct must has tag "db",what map to field of database.
func (build *UpdateSqlBuilder) UpdateByStruct(tableMap interface{}, skipEmpty bool) {
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
			log.Printf("%s has no tag 'db', skip it.\n", tableType.Field(i).Name)
			continue
		}
		value := tableValue.Field(i).Interface()
		if skipEmpty == true && (value == 0 || value == "") {
			continue
		}
		if sqlStr == "" {
			sqlStr = fmt.Sprintf("`%s`=%s", dbTag, build.flag)
		} else {
			sqlStr = fmt.Sprintf("%s,`%s`=%s", sqlStr, dbTag, build.flag)
		}
		build.args = append(build.args, value)
	}
	build.handleStr = sqlStr
}

// Build sql string of update with one condition
func (build *UpdateSqlBuilder) UpdateSet(fieldName string, fieldValue interface{}) {
	if build.sqlType != SQL_TYPE_UPDATE {
		log.Fatalf("Builder is for %s, can not use Update!", build.sqlType)
	}

	if build.handleStr == "" {
		build.handleStr = fmt.Sprintf("%s=%s", fieldName, build.flag)
	} else {
		build.handleStr = fmt.Sprintf("%s,%s=%s", build.handleStr, fieldName, build.flag)
	}
	build.args = append(build.args, fieldValue)
}

// Build sql string of insert by struct
// The value of struct must has tag "db",what map to field of database.
func (build *InsertSqlBuilder) InsertByStruct(tableMap interface{}) {
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
			log.Printf("%s has no tag 'db', skip it.\n", tableType.Field(i).Name)
			continue
		}
		if sqlStr == "" {
			sqlStr = fmt.Sprintf("`%s`", dbTag)
			valStr = build.flag
		} else {
			sqlStr = fmt.Sprintf("%s,`%s`", sqlStr, dbTag)
			valStr = fmt.Sprintf("%s,%s", valStr, build.flag)
		}
		build.args = append(build.args, tableValue.Field(i).Interface())
	}
	build.handleStr = fmt.Sprintf("INSERT INTO `%s`(%s) VALUES(%s);", build.table, sqlStr, valStr)
}

func (build *SelectSqlBuilder) SelectByStruct(tableMap interface{}, skipEmpty bool) {
	if build.sqlType != SQL_TYPE_SELECT {
		log.Fatalf("SQL type error")
	}

	tableType := reflect.TypeOf(tableMap)
	valueType := reflect.ValueOf(tableMap)

	num := tableType.NumField()

	var fieldStr string
	for i := 0; i < num; i++ {
		dbTag := tableType.Field(i).Tag.Get("db")
		if dbTag == "" {
			log.Printf("%s has no tag 'db', skip it.\n", tableType.Field(i).Name)
			continue
		}
		if fieldStr == "" {
			fieldStr = fmt.Sprintf("`%s`", dbTag)
		} else {
			fieldStr = fmt.Sprintf("%s,`%s`", fieldStr, dbTag)
		}
		value := valueType.Field(i).Interface()
		if skipEmpty == true && (value == "" || value == 0) {
			continue
		} else {
			build.WhereEq(dbTag, value)
		}
	}
	build.fields = fieldStr
}

// Build sql string of 'where' with '='
func (build *whereBuilder) WhereEq(fieldName string, fieldValue interface{}) {
	build.buildWhereCondition(fieldName, WHERE_TYPE_EQ, fieldValue)
}

// Build sql string of 'where' with '<>'
func (build *whereBuilder) WhereNeq(fieldName string, fieldValue interface{}) {
	build.buildWhereCondition(fieldName, WHERE_TYPE_NEQ, fieldValue)
}

// Build sql string of 'where' with '>'
func (build *whereBuilder) WhereGt(fieldName string, fieldValue interface{}) {
	build.buildWhereCondition(fieldName, WHERE_TYPE_GT, fieldValue)
}

// Build sql string of 'where' with '>='
func (build *whereBuilder) WhereGte(fieldName string, fieldValue interface{}) {
	build.buildWhereCondition(fieldName, WHERE_TYPE_GTE, fieldValue)
}

// Build sql string of 'where' with '<'
func (build *whereBuilder) WhereLt(fieldName string, fieldValue interface{}) {
	build.buildWhereCondition(fieldName, WHERE_TYPE_LT, fieldValue)
}

// Build sql string of 'where' with '<='
func (build *whereBuilder) WhereLte(fieldName string, fieldValue interface{}) {
	build.buildWhereCondition(fieldName, WHERE_TYPE_LTE, fieldValue)
}

// Build sql string of 'where' with 'like'
func (build *whereBuilder) WhereLike(fieldName string, fieldValue interface{}) {
	build.buildWhereCondition(fieldName, WHERE_TYPE_Like, fieldValue)
}

// Build sql string of 'where' with 'like'
func (build *whereBuilder) WhereLikeBefore(fieldName string, fieldValue interface{}) {
	build.buildWhereCondition(fieldName, WHERE_TYPE_LikeBEF, fieldValue)
}

// Build sql string of 'where' with 'like'
func (build *whereBuilder) WhereLikeAfter(fieldName string, fieldValue interface{}) {
	build.buildWhereCondition(fieldName, WHERE_TYPE_LikeAFT, fieldValue)
}

// Build sql string of 'in' condition with conditions
func (build *whereBuilder) WhereIn(fieldName string, args []interface{}) {
	length := len(args)
	if length < 1 {
		log.Fatalf("Need args")
	}

	condition := ""
	for i := 0; i < length; i++ {
		if args[i] == "" {
			continue
		}
		if condition == "" {
			condition = build.flag
		} else {
			condition = fmt.Sprintf("%s,%s", condition, build.flag)
		}
		build.args = append(build.args, args[i])
	}

	if build.whereStr == "" {
		build.whereStr = fmt.Sprintf("%s IN (%s)", fieldName, condition)
	} else {
		build.whereStr = fmt.Sprintf("%s AND %s IN (%s)", build.whereStr, fieldName, condition)
	}
}

// Build sql string of 'or' condition with struct WhereOrCondition
func (build *whereBuilder) WhereOr(args []WhereOrCondition) {
	var orStr string
	for _, value := range args {
		condition := getWhereTypeString(value.WhereType)
		fieldValue := value.FieldValue
		switch value.WhereType {
		case WHERE_TYPE_Like:
			fieldValue = fmt.Sprintf("%%%s%%", fieldValue)
			break
		case WHERE_TYPE_LikeAFT:
			fieldValue = fmt.Sprintf("%s%%", fieldValue)
			break
		case WHERE_TYPE_LikeBEF:
			fieldValue = fmt.Sprintf("%%%s", fieldValue)
			break
		default:
			break
		}

		conditionStr := fmt.Sprintf("`%s`%s%s", value.FieldName, condition, build.flag)
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
		build.whereStr = fmt.Sprintf("(%s) OR (%s)", build.whereStr, orStr)
	}
}

// Build sql string with struct, whick has tag "db"
func (build *whereBuilder) WhereByStruct(tableMap interface{}, skipEmpty bool) {
	tableType := reflect.TypeOf(tableMap)
	valueType := reflect.ValueOf(tableMap)

	num := tableType.NumField()

	for i := 0; i < num; i++ {
		dbTag := tableType.Field(i).Tag.Get("db")
		if dbTag == "" {
			log.Printf("%s has no tag 'db', skip\n", tableType.Field(i).Name)
			continue
		}
		value := valueType.Field(i).Interface()
		if skipEmpty == true && (value == "" || value == 0) {
			continue
		}
		build.WhereEq(dbTag, value)
	}
}

// Get sql string
func (build *sqlBuilder) String() string {
	whereStr := ""
	if build.whereStr != "" {
		whereStr = fmt.Sprintf(" WHERE %s", build.whereStr)
	}
	switch build.sqlType {
	case SQL_TYPE_INSERT:
		return build.handleStr
	case SQL_TYPE_UPDATE:
		return fmt.Sprintf("UPDATE `%s` SET %s%s;", build.table, build.handleStr, whereStr)
	case SQL_TYPE_SELECT:
		return fmt.Sprintf("SELECT %s FROM `%s`%s;", build.fields, build.table, whereStr)
	case SQL_TYPE_DELETE:
		return fmt.Sprintf("DELETE FROM `%s`%s;", build.table, whereStr)
	}
	return ""
}

// Get all Args
func (build *sqlBuilder) Args() []interface{} {
	return build.args
}

// Set sql flag
func (build *sqlBuilder) SetFlag(flag string) {
	build.flag = flag
}

// Set search fields
func (build *SelectSqlBuilder) SetSearchFields(selectField string) {
	build.fields = selectField
}

func (build *SelectSqlBuilder) Limit(offset int64, size int64) {
	build.whereStr = fmt.Sprintf("%s LIMIT %d,%d", build.whereStr, offset, size)
}

func (build *SelectSqlBuilder) GroupBy(fieldName string) {
	build.whereStr = fmt.Sprintf("%s GROUP BY %s", build.whereStr, fieldName)
}

func (build *SelectSqlBuilder) OrderBy(orderBy string) {
	build.whereStr = fmt.Sprintf("%s ORDER BY %s", build.whereStr, orderBy)
}

// Build condition string
func (build *whereBuilder) buildWhereCondition(fieldName string, whereType WhereType, fieldValue interface{}) {
	if whereType == WHERE_TYPE {
		log.Fatalf("Where type error")
	}

	condition := getWhereTypeString(whereType)
	switch whereType {
	case WHERE_TYPE_Like:
		fieldValue = fmt.Sprintf("%%%s%%", fieldValue)
		break
	case WHERE_TYPE_LikeAFT:
		fieldValue = fmt.Sprintf("%s%%", fieldValue)
		break
	case WHERE_TYPE_LikeBEF:
		fieldValue = fmt.Sprintf("%%%s", fieldValue)
		break
	default:
		break
	}

	conditionStr := fmt.Sprintf("`%s`%s%s", fieldName, condition, build.flag)
	if build.whereStr == "" {
		build.whereStr = conditionStr
	} else {
		build.whereStr = fmt.Sprintf("%s AND %s", build.whereStr, conditionStr)
	}
	build.args = append(build.args, fieldValue)
}

func getWhereTypeString(whereType WhereType) string {
	var condition string
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
		condition = " LIKE "
		break
	case WHERE_TYPE_LikeAFT:
		condition = " LIKE "
		break
	case WHERE_TYPE_LikeBEF:
		condition = " LIKE "
		break
	default:
		condition = "="
		break
	}

	return condition
}
