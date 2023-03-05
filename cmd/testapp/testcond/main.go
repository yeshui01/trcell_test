/**
 * @author: [mknight]
 * @email : [824338670@qq.com]
 * @create:	2023-01-18 09:34:55
 * @modify:	2023-01-18 09:34:55
 * @desc  :	[test use condition var]
 */
package main

import (
	"fmt"
	"sync"
	"time"
)

func testFor() {
	num := 0
	for num != 1 {
		num = num + 2
	}

	fmt.Print("num is ", num)
}

func main() {
	// testFor()
	fmt.Printf("hello sync and condition var\n")
	c := sync.NewCond(&sync.Mutex{})
	var num int
	var wg sync.WaitGroup
	for i := 1; i <= 2; i++ {
		wg.Add(1)
		go func(id int) {
			fmt.Println("Enter Thread ID:", id)
			c.L.Lock()
			for num == 0 {
				fmt.Println("Enter loop: Thread ID:", id, "num:", num)
				c.Wait()
				fmt.Println("Exit loop: Thread ID:", id, "num:", num)
			}
			num++
			c.L.Unlock()
			fmt.Println("Exit Thread ID:", id, ", num:", num)
			wg.Done()
		}(i)
	}

	time.Sleep(time.Second * 2)
	fmt.Println("Sleep 2 second")
	c.L.Lock()
	num++
	c.Broadcast()
	c.L.Unlock()
	time.Sleep(time.Second)
	wg.Wait()
	fmt.Println("Program normal exit")
}
