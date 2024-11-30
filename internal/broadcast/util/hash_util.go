package util

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

// 计算消息结构的哈希值
func CalculateHash(message interface{}) (string, error) {
	// 将消息结构序列化为 JSON
	data, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	// 计算 SHA-256 哈希值
	hash := sha256.Sum256(data)

	// 返回哈希值的十六进制字符串
	return fmt.Sprintf("%x", hash), nil
}
