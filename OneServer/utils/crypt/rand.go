package crypt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func GenerateUID(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	// 获取当前时间戳（纳秒级）
	timestamp := time.Now().UnixNano()

	// 将时间戳转换为16进制字符串
	timestampHex := fmt.Sprintf("%x", timestamp)

	// 计算还需要多少随机字符
	randomLength := length - len(timestampHex)

	// 如果时间戳部分已经超过了请求的长度，则截断
	if randomLength <= 0 {
		return timestampHex[:length], nil
	}

	// 计算需要多少随机字节
	byteLength := randomLength / 2
	if randomLength%2 != 0 {
		byteLength++
	}

	// 生成随机字节
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// 转换为十六进制字符串
	randomHex := hex.EncodeToString(randomBytes)[:randomLength]

	// 拼接时间戳和随机部分
	result := timestampHex + randomHex

	// 确保长度正确
	if len(result) > length {
		result = result[:length]
	}

	return result, nil
}
