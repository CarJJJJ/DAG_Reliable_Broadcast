package brachabroadcast

import (
	"log"
	"time"
)

func (instance *NodeExtention) CountTPS() {
	// 每秒统计量级
	for {
		log.Printf("[INFO] tps :%v, count:%v, lastSecond count:%v",
			instance.ReliableBroadcastCount-instance.ReliableBroadcastCountLastSecond,
			instance.ReliableBroadcastCount,
			instance.ReliableBroadcastCountLastSecond,
		)
		instance.ReliableBroadcastCountLastSecond = instance.ReliableBroadcastCount
		time.Sleep(time.Second)
	}
}
