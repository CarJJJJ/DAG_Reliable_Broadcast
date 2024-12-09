package test_util

import (
	"testing"

	"github.com/klauspost/reedsolomon"
)

func TestReedSolomon(t *testing.T) {
	// n-t、t
	encoder, err := reedsolomon.New(3, 1)
	if err != nil {
		t.Fatal(err)
	}

	// 确保所有数据分片长度相同（使用填充）
	data := make([][]byte, 4)
	data[0] = []byte("hello ") // 添加空格使长度为6
	data[1] = []byte("world ") // 添加空格使长度为6
	data[2] = []byte("foo   ") // 添加空格使长度为6
	data[3] = make([]byte, 6)  // 为校验分片预分配相同大小的空间

	err = encoder.Encode(data)
	if err != nil {
		t.Fatal(err)
	}

	data[0] = nil

	// 打印data
	for _, shard := range data {
		t.Logf("shard: %s", string(shard))
	}

	err = encoder.Reconstruct(data)
	if err != nil {
		t.Fatal(err)
	}

	// 打印data
	for _, shard := range data {
		t.Logf("shard: %s", string(shard))
	}

}
