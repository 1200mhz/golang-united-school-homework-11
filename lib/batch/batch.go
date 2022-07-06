package batch

import (
	"context"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type users struct {
	items []user
	mu    sync.Mutex
}

type user struct {
	ID int64
}

func getOne(id int64) user {
	time.Sleep(time.Millisecond * 100)
	return user{ID: id}
}

func getBatch(n int64, pool int64) (res []user) {
	result := &users{}

	eg, _ := errgroup.WithContext(context.Background())
	eg.SetLimit(int(pool))

	var i int64
	for i = 0; i < n; i++ {
		func(i int64) {
			eg.Go(func() error {
				u := getOne(i)
				result.mu.Lock()
				result.items = append(result.items, u)
				result.mu.Unlock()

				return nil
			})
		}(i)
	}

	eg.Wait()

	return result.items
}

func getBatch2(n int64, pool int64) (res []user) {
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
