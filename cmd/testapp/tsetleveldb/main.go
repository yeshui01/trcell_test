package main

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	_ "github.com/syndtr/goleveldb/leveldb/cache"
	lvldbfilter "github.com/syndtr/goleveldb/leveldb/filter"
	lvldbopt "github.com/syndtr/goleveldb/leveldb/opt"
)

func main() {
	fmt.Println("hello leveldb")
	optSetting := &lvldbopt.Options{
		Filter: lvldbfilter.NewBloomFilter(8), // key过滤策略,这里设置8个big,减少磁盘扫描次数,提高效率
		// OpenFilesCacher: lvldbcache.NewLRU(100 * 1048576), // 100 MB缓存
		OpenFilesCacheCapacity: 100 * 1048576, // 设置100 MB缓存
	}
	db, err := leveldb.OpenFile("/root/self_projects/trcell/lvldb", optSetting)
	if err != nil {
		return
	}
	// 具体操作示例可以参考github.com/syndtr/goleveldb/leveldb中的README.md
	// PUT

	// Get

	// delete

	// BatchWrite

	// iterator

	defer db.Close()
}
