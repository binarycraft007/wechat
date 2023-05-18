package utils

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func GetDeviceID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("e%.15s", fmt.Sprintf("%0.15f", rand.Float64())[2:17])
}

func GetClientMsgId() int64 {
	milliseconds := time.Now().UnixNano() / int64(time.Millisecond)
	return int64(float64(milliseconds) * 1e3)
}

func GetErrorMsgInt(code int) string {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d, %s error: %d", file, line, f.Name(), code)
}

func GetErrorMsgStr(str string) string {
	pc, file, line, _ := runtime.Caller(1)
	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%s:%d, %s error: %s", file, line, f.Name(), str)
}
