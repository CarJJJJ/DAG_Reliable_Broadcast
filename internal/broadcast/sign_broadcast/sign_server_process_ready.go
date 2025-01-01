package signbroadcast

func (node *NodeExtention) ProcessReadyMessage() {

	readyMessage := <-node.BCBReadyPool
	// log.Printf("[INFO] 收到ReadyMessage, unique_index: %s, from: %d", readyMessage.UniqueIndex, readyMessage.NodeID)

	// node.RecvReadyMessageForUniqueIndexNumberSync.Upsert(readyMessage.UniqueIndex, cmap.New[int](), func(exist bool, valueInMap, newValue cmap.ConcurrentMap[string, int]) cmap.ConcurrentMap[string, int] {
	// 	if exist {
	// 		// 如果存在，直接使用现有的 map
	// 		valueInMap.Set(strconv.Itoa(readyMessage.NodeID), 1)
	// 		return valueInMap
	// 	} else {
	// 		// 如果不存在，使用新的 map
	// 		newValue.Set(strconv.Itoa(readyMessage.NodeID), 1)
	// 		return newValue
	// 	}
	// })

	node.RecvReadyMessageForUniqueIndexNumberMu.Lock()
	if _, ok := node.RecvReadyMessageForUniqueIndexNumber[readyMessage.UniqueIndex]; !ok {
		node.RecvReadyMessageForUniqueIndexNumber[readyMessage.UniqueIndex] = make(map[int]int)
		node.RecvReadyMessageForUniqueIndexNumber[readyMessage.UniqueIndex][readyMessage.NodeID] = 1
	} else {
		node.RecvReadyMessageForUniqueIndexNumber[readyMessage.UniqueIndex][readyMessage.NodeID] = 1
	}
	node.RecvReadyMessageForUniqueIndexNumberMu.Unlock()

	// 如果收到n-t个ready，那么就将消息宣布可靠广播
	if number := len(node.RecvReadyMessageForUniqueIndexNumber[readyMessage.UniqueIndex]); number >= node.N-node.T {
		// 如果还没有宣布过可靠广播
		if _, ok := node.HadReliableBroadcastUniqueIndex.Get(readyMessage.UniqueIndex); !ok {
			// log.Printf("[INFO] 收到n-t个ready, 宣布可靠广播, unique_index: %s", readyMessage.UniqueIndex)
			// 记录已经宣布过可靠广播
			node.HadReliableBroadcastUniqueIndex.Set(readyMessage.UniqueIndex, 1)
			node.ReliableBroadcastCount++
		}
	}

}
