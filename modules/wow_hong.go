package modules

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Config struct{
	Host string `json:"host"`
	Port string `json:"port"`
	UserName string `json:"userName"`
	Pwd string `json:"pwd"`
}

// 访问日志
type LoginLog struct {
	ID int32 `json:"id"`
	// ip地址
	IP string `json:"ip"`
	// 访问主体接口
	Method string `json:"method"`
	// 创建时间
	CreateTime timestamp.Timestamp `json:"createTime"`
}

// 统一处理错误
func CheckError(c *gin.Context, err error) {
	if err != nil {
		log.Fatalf("error: %s", err.Error())
		c.String(http.StatusForbidden, err.Error())
	}
}
