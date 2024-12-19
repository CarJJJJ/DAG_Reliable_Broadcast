package dagbroadcast

import (
	"time"
)

func (instance *NodeExtention) CountTPS() {
	// 每秒统计量级
	for {
		time.Sleep(time.Second)
	}
}
