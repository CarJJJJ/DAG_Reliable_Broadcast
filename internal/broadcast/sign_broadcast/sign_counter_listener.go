package signbroadcast

import (
	"log"
	"time"
)

func (Instance *NodeExtention) CountTPS() {
	// 每秒统计量级
	for {
		log.Printf("[INFO] tps :%v, count:%v, lastSecond count:%v",
			Instance.ReliableBroadcastCount-Instance.ReliableBroadcastCountLastSecond,
			Instance.ReliableBroadcastCount,
			Instance.ReliableBroadcastCountLastSecond,
		)
		Instance.ReliableBroadcastCountLastSecond = Instance.ReliableBroadcastCount
		time.Sleep(time.Second)
	}
}
