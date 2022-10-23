/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-16 14:18:14
 * @LastEditTime: 2022-09-20 17:59:40
 * @FilePath: \trcell\pkg\utils\string_to_hash_num.go
 */
package utils

import (
	"encoding/binary"
	"hash/fnv"
)

func StrToHashInt(originStr string) (uint16, uint16) {
	if len(originStr) < 1 {
		return 0, 0
	}
	byteSeque := []byte(originStr)
	if len(byteSeque) < 2 {
		return 0, uint16(uint8(byteSeque[0]))
	}
	if len(byteSeque) < 3 {
		return 0, binary.LittleEndian.Uint16(byteSeque[0:])
	}
	if len(byteSeque) < 4 {
		return uint16(uint8(byteSeque[0])), binary.LittleEndian.Uint16(byteSeque[1:])
	}

	high := binary.LittleEndian.Uint16(byteSeque[0:])
	low := binary.LittleEndian.Uint16(byteSeque[2:])
	return high, low
	// return uint32(high)<<16 + uint32(low)
}

// FNV hash
func Fnv32Value(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func StringHashToInt32(srcData string) uint32 {
	hash32 := fnv.New32()
	hash32.Write([]byte(srcData))
	return hash32.Sum32()
}
