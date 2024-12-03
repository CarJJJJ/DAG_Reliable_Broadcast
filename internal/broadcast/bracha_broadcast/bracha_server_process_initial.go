package brachabroadcast

import (
	"DAG_Reliable_Broadcast/internal/broadcast/util"
	"log"
)

func (instance *NodeExtention) ProcessInitial() {
	select {
	case msg := <-instance.InitialPool:
		// log.Printf("[INFO] 收到初始消息: %+v", msg) // 记录收到的初始消息

		// 计算哈希值
		hash, err := util.CalculateHash(msg)
		if err != nil {
			log.Printf("[ERROR] 计算哈希值失败: %v", err) // 记录错误
			return
		}

		// 检查哈希值是否已经存在于 HadEchoInitial 中
		if _, exists := instance.HadEchoInitial.Get(hash); exists {
			// log.Printf("[INFO] 消息:%v,已处理,不再处理", msg) // 如果哈希值存在，记录日志
			return
		}

		// 如果哈希值不存在，放入 map 中并处理消息
		instance.HadEchoInitial.Set(hash, 1)

		// 封装Echo
		echo := EchoMessage{
			Type:           echo_type,
			InitialMessage: msg,
			NodeID:         instance.Node.Id,
		}

		// 广播echo
		instance.BroadcastEchoToServers(echo)

	default:
		log.Println("[INFO] 当前没有初始消息可处理") // 如果没有消息，记录日志
	}
}
