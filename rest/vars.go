package rest

import (
	"os"
	"strconv"
	"time"

	"github.com/fimreal/goutils/ezap"
)

var (
	ExpTime     = 2 * time.Minute
	RefreshTime = 14 * 24 * time.Hour
	Secret      string
)

func init() {
	Secret = os.Getenv("SECRET")
	if Secret == "" {
		Secret = strconv.Itoa(time.Now().Nanosecond())
		os.Setenv("SECRET", Secret)
		ezap.SetLogTime("2006-01-02 15:04:05")
		ezap.Warn("未配置 SECRET 变量值, 自动设置为 ", Secret)
	}
}
