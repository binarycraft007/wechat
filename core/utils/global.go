package utils

import (
	"fmt"
	"math/rand"
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
