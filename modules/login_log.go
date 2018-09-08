package modules

import "github.com/illidan33/sql-builder"

type LoginLog struct {
	// 访问IP
	IP string `json:"ip" db:"ip"`
	// 访问页面
	Method string `json:"method" db:"method"`
	// 访问时间
	CreateTime string `json:"createTime" db:"createtime"`
}

// 创建登录日志
func CreateLog(log LoginLog) bool {
	var build sql_builder.SqlBuilder
	build.Init("login_log", sql_builder.SQL_TYPE_INSERT)
	build.InsertByStruct(log)

	res, err := DbConn.Exec(build.String(), build.Args()...)
	rowNum, err := res.RowsAffected()
	if err != nil {
		CheckErr(build.String(), err)
		return false
	}

	if rowNum > 0 {
		return true
	}

	return false
}
