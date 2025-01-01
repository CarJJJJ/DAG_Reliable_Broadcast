package signbroadcast

func (node *NodeExtention) ProcessReconstruct() {

	msg := <-node.BCBReconstructPool
	// log.Printf("[INFO] 收到重构消息, unique_index: %s, from: %d", msg.UniqueIndex, msg.NodeID)

	// 把分片数据放到map，等待恢复数据

	// node.ReconstructDataForUniqueIndexSync.Upsert(msg.UniqueIndex, make([][]byte, node.N), func(exist bool, valueInMap, newValue [][]byte) [][]byte {
	// 	if exist {
	// 		// 如果存在，直接使用现有的 slice
	// 		valueInMap[msg.NodeID] = msg.DataFrom
	// 		return valueInMap
	// 	} else {
	// 		// 如果不存在，使用新的 slice
	// 		newValue[msg.NodeID] = msg.DataFrom
	// 		return newValue
	// 	}
	// })

	node.ReconstructDataForUniqueIndexMu.Lock()
	if _, ok := node.ReconstructDataForUniqueIndex[msg.UniqueIndex]; !ok {
		node.ReconstructDataForUniqueIndex[msg.UniqueIndex] = make([][]byte, node.N)
		node.ReconstructDataForUniqueIndex[msg.UniqueIndex][msg.NodeID] = msg.DataFrom
	} else {
		node.ReconstructDataForUniqueIndex[msg.UniqueIndex][msg.NodeID] = msg.DataFrom
	}
	node.ReconstructDataForUniqueIndexMu.Unlock()

	// 记录收到分片消息的数量
	if number, ok := node.RecvReconstructMessageForUniqueIndexNumber.Get(msg.UniqueIndex); ok {
		node.RecvReconstructMessageForUniqueIndexNumber.Set(msg.UniqueIndex, number+1)
	} else {
		node.RecvReconstructMessageForUniqueIndexNumber.Set(msg.UniqueIndex, 1)
	}

	// 当收到N-T份数据的时候，就可以恢复数据了
	if number, ok := node.RecvReconstructMessageForUniqueIndexNumber.Get(msg.UniqueIndex); ok && number >= node.N-node.T {
		// 	恢复数据
		data := node.ReconstructDataForUniqueIndex[msg.UniqueIndex]
		node.ReedSolomonEncoder.Reconstruct(data)

		// 广播Ready及恢复后的数据
		readyMessage := ReadyMessage{
			Type:        BCBReadyType,
			UniqueIndex: msg.UniqueIndex,
			Message:     data,
			NodeID:      node.Node.Id,
		}
		// 发的时候再检查一遍, 如果还没有Ready过
		if _, ok := node.HadReadyUniqueIndex.Get(msg.UniqueIndex); !ok {
			// log.Printf("[INFO],发送Ready, unique_index: %s,", readyMessage.UniqueIndex)
			node.BroadcastReadyMessage(readyMessage)
			// 记录已经Ready
			node.HadReadyUniqueIndex.Set(msg.UniqueIndex, 1)
		}

	}

	// // 如果还没有Disperse，则编码、发送
	// if _, ok := node.HadDisperseUniqueIndex.Get(msg.UniqueIndex); !ok {
	// 	// 打印amplifier的日志
	// 	log.Printf("[INFO] Reconstruct后还没有Disperse, 所以编码、发送, unique_index: %s", msg.UniqueIndex)
	// 	// 编码
	// 	err := node.ReedSolomonEncoder.Encode(node.ReconstructDataForUniqueIndex[msg.UniqueIndex])
	// 	if err != nil {
	// 		log.Printf("[ERROR] 编码失败, unique_index: %s, 错误: %v", msg.UniqueIndex, err)
	// 		return
	// 	}
	// 	// 发送
	// 	for i := 0; i < node.N; i++ {
	// 		disperseMessage := DisperseMessage{
	// 			UniqueIndex: msg.UniqueIndex,
	// 			Type:        BCBDisperseType,
	// 			DataFrom:    node.ReconstructDataForUniqueIndex[msg.UniqueIndex][i],
	// 			NodeID:      msg.NodeID,
	// 		}
	// 		node.SendDisperseMessage(disperseMessage, i)
	// 	}

	// 	// 记录已经Disperse
	// 	node.HadDisperseUniqueIndex.Set(msg.UniqueIndex, 1)
	// }

}
