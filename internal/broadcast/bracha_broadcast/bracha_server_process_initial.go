package brachabroadcast

import "log"

func (instance *NodeExtention) ProcessInitial() {
	select {
	case msg := <-instance.InitialPool:
		log.Printf("[INFO] 收到初始消息: %+v", msg) // 记录收到的初始消息
	default:
		log.Println("[INFO] 当前没有初始消息可处理") // 如果没有消息，记录日志
	}
}
