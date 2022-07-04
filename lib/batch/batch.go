package batch

import (
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

func getBatch(n int64, pool int64) (res []user) {
	type users struct {
		items []user
		mu    sync.Mutex
	}

	result := &users{}
	var waitGroup sync.WaitGroup
	ch := make(chan struct{}, pool)

	var i int64
	for i = 0; i < n; i++ {
		waitGroup.Add(1)

		ch <- struct{}{}

		go func(i int64) {
			u := getOne(i)

			<-ch

			result.mu.Lock()
			result.items = append(result.items, u)
			result.mu.Unlock()

			waitGroup.Done()
		}(i)
	}

	waitGroup.Wait()

	return result.items
}
