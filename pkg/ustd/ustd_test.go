/*
 * @Author: mknight(tianyh)
 * @Mail: 824338670@qq.com
 * @Date: 2022-09-20 17:19:03
 * @LastEditTime: 2022-09-26 15:06:09
 * @FilePath: \trcell\pkg\ustd\ustd_test.go
 */
// /*
//  * @Author: mknight(tianyh)
//  * @Mail: 824338670@qq.com
//  * @Date: 2022-09-20 16:14:38
//  * @LastEditTime: 2022-09-20 18:28:20
//  * @FilePath: \trcell\pkg\ustd\ustd_test.go
//  */
package ustd

// import (
// 	"fmt"
// 	"hash/fnv"
// 	"testing"
// 	xstd "trcell/pkg/ustd"
// )

// func TestVectorTest(t *testing.T) {
// 	fmt.Println("TestVectorTest")
// 	stdVec := xstd.NewVector[int32]()
// 	stdVec.PushBack(1)
// 	stdVec.PushBack(2)
// 	stdVec.ForEach(func(idx int32, v int32) {
// 		fmt.Printf("idx:%d, value:%d\n", idx, v)
// 	})
// }
// func TestMapT(t *testing.T) {
// 	fmt.Println("TestMaMapT")
// 	stdMaMapT := xstd.NewMap[int32, string]()
// 	stdMaMapT.Set(1, "value1")
// 	stdMaMapT.Set(2, "value2")
// 	stdMaMapT.ForEach(func(k int32, v string) {
// 		fmt.Printf("key:%d, value:%s\n", k, v)
// 	})
// }
// func Fnv32Value(key string) uint32 {
// 	hash := uint32(2166136261)
// 	const prime32 = uint32(16777619)
// 	for i := 0; i < len(key); i++ {
// 		hash *= prime32
// 		hash ^= uint32(key[i])
// 	}
// 	return hash
// }

// func StringHashToInt32(srcData string) uint32 {
// 	hash32 := fnv.New32()
// 	hash32.Write([]byte(srcData))
// 	return hash32.Sum32()
// }

// func TestString32(t *testing.T) {
// 	// fmt.Printf("abcdef to 32:%d\n", Fnv32Value("abcdef"))
// 	// fmt.Printf("ghikj to 32:%d\n", Fnv32Value("ghikj"))
// 	fmt.Printf("abcdef to 32:%d\n", StringHashToInt32("abcdef"))
// 	fmt.Printf("ghikj to 32:%d\n", StringHashToInt32("ghikj"))
// }
