package miner

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func Miner(
	ctx context.Context,
	wg *sync.WaitGroup,
	transferPoint chan<- int,
	id int,
	coalPerMining int,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("I'm miner", id, "my workday is over")
			return
		default:
			fmt.Println("I'm miner", id, "I'm starting to mine coal!")
			time.Sleep(1 * time.Second)
			fmt.Println("I'm miner", id, "Coal mined: ", coalPerMining)

			transferPoint <- coalPerMining
			fmt.Println("I'm miner", id, "Coal delivered: ", coalPerMining)
		}
	}
}

func MinerPool(ctx context.Context, minerCount int) <-chan int {
	coalTransferPoint := make(chan int)

	wg := &sync.WaitGroup{}

	for i := 1; i <= minerCount; i++ {
		wg.Add(1)
		go Miner(ctx, wg, coalTransferPoint, i, i*10)
	}

	go func() {
		wg.Wait()
		close(coalTransferPoint)
	}()

	return coalTransferPoint
}
