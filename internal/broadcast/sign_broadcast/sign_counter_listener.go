package signbroadcast

import (
	"time"
)

func (Instance *NodeExtention) CountTPS() {
	// 每秒统计量级
	for {
		time.Sleep(time.Second)
	}
}
