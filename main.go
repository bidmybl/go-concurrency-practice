package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"concurrency/miner"
	"concurrency/postman"
)

func main() {
	var coal atomic.Int64
	var mails []string
	mtx := sync.Mutex{}

	minerContext, minerCancel := context.WithTimeout(context.Background(), 3*time.Second)
	postmanContext, postmanCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer minerCancel()
	defer postmanCancel()

	coalTransferPoint := miner.MinerPool(minerContext, 20)
	mailTransferPoint := postman.PostmanPool(postmanContext, 20)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for c := range coalTransferPoint {
			coal.Add(int64(c))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for m := range mailTransferPoint {
			mtx.Lock()
			mails = append(mails, m)
			mtx.Unlock()
		}
	}()

	wg.Wait()

	fmt.Println("Coal: ", coal.Load())

	mtx.Lock()
	fmt.Println("Mails: ", len(mails))
	mtx.Unlock()
}
