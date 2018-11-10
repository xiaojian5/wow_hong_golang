package modules

import (
	"github.com/illidan33/sql-builder"
	"time"
)

type LoginForGet struct {
	ID        int    `json:"id" db:"id"`
	IP        string `json:"ip" db:"ip"`
	Method    string `json:"method" db:"method"`
	LoginDate string `json:"loginDate" db:"login_date"`
	Count     int    `json:"count" db:"count"`
}

type LoginLog struct {
	ID         int    `json:"id" db:"id"`
	IP         string `json:"ip" db:"ip"`
	Method     string `json:"method" db:"method"`
	LoginDate  string `json:"loginDate" db:"login_date"`
	Count      int    `json:"count" db:"count"`
	CreateTime string `json:"createTime" db:"create_time"`
	UpdateTime string `json:"updateTime" db:"update_time"`
}

var(
	table = "macro_login_log"
)

// 创建登录日志
func CreateLog(ip string, method string) error {
	now := time.Now()
	t := now.Format("2006-01-02 15:04:05")

	log := LoginLog{
		ID:         0,
		IP:         ip,
		Method:     method,
		LoginDate:  now.Format("2006-01-02"),
		Count:      1,
		CreateTime: t,
		UpdateTime: t,
	}
	builder := sql_builder.Insert(table)
	builder.InsertByStruct(log)

	conn := GetDbConn()
	_, err := conn.Exec(builder.String(), builder.Args()...)
	if err != nil {
		return err
	}

	return nil
}

func UpdateLog(id int, count int) error {
	builder := sql_builder.Update(table)
	builder.WhereEq("id", id)
	builder.UpdateSet("count", count)

	conn := GetDbConn()
	_, err := conn.Exec(builder.String(), builder.Args()...)
	if err != nil {
		return err
	}

	return nil
}

func GetLog(ip string, method string, date string) (LoginForGet, error) {
	builder := sql_builder.Select(table)
	builder.WhereEq("ip", ip)
	builder.WhereEq("method", method)
	builder.WhereEq("login_date", date)
	builder.SetSearchFields([]string{"id", "ip", "method", "login_date", "count"})

	log := LoginForGet{}
	conn := GetDbConn()
	err := conn.Get(&log, builder.String(), builder.Args()...)
	if err != nil {
		return LoginForGet{}, err
	}
	return log, nil
}