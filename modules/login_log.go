package modules

type LoginLob struct {
	// 访问IP
	IP string `json:"ip" db:"ip"`
	// 访问页面
	Method string `json:"method" db:"method"`
	// 访问时间
	CreateTime string `json:"createTime" db:"createtime"`
}

// 创建登录日志
func CreateLog(log LoginLob) bool {
	sqlStr := "INSERT INTO login_log(`ip`,`method`,`createtime`) VALUES(:IP,:Method,:CreateTime)"
	res, err := DbConn.NamedExec(sqlStr, log)
	rowNum, err := res.RowsAffected()
	if err != nil {
		CheckErr(sqlStr, err)
		return false
	}

	if rowNum > 0 {
		return true
	}

	return false
}
