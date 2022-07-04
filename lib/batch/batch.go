package batch

import (
	"fmt"
	"sync"
	"time"
)

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

/*func main() {
	getBatch(10, 2)
}*/

func getBatch(n int64, pool int64) (res []user) {
	type users struct {
		items []user
		mu    sync.Mutex
	}

	result := &users{}
	var waitGroup sync.WaitGroup
	ch := make(chan int64, pool)

	go func() {
		for {
			select {
			case i := <-ch:
				u := getOne(i)

				result.mu.Lock()
				result.items = append(result.items, u)
				result.mu.Unlock()

				waitGroup.Done()
			}
		}
	}()

	var i int64
	for i = 0; i < n; i++ {
		waitGroup.Add(1)
		go func(i int64) {
			ch <- i
		}(i)
	}

	waitGroup.Wait()

	fmt.Println("RES:", result.items)

	return result.items
}
