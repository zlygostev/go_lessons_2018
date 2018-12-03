package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

func periodicJob1(ctx context.Context, waitGroup *sync.WaitGroup, log io.Writer, threadID int, secs int) {
	defer waitGroup.Done()
	fmt.Fprintf(log, "[%d] Periodic Job is Started...\n", threadID)
	//	ticker := time.Tick(time.Duration(secs) * time.Second)
	ticker := time.NewTicker(time.Duration(secs) * time.Second)
	counter := 0
LOOP:
	for {
		select {
		//		case tickerTime := <-ticker:
		case tickerTime := <-ticker.C:
			fmt.Fprintf(log, "[%d] Step %d at %v...\n", threadID, counter, tickerTime)
			counter++
		case <-ctx.Done():
			// ticker.Stop()
			fmt.Fprintf(log, "[%d] Stop signal come\n", threadID)
			break LOOP
		}
	}

	fmt.Fprintf(log, "[%d] Periodic Job has been finished...\n", threadID)
}

func main() {
	out := os.Stdout
	wg := &sync.WaitGroup{}
	thID := 1
	ctx, finishFunc := context.WithCancel(context.Background())
	wg.Add(1)
	go periodicJob1(ctx, wg, out, thID, 2)
	// for thID := 2; thID < threadsCount; {
	// 	go doJob(out, thID, msgChannel)
	// 	thID++
	// }
	// afterFuncTimer := time.AfterFunc(1*time.Second, func() {
	// 	afterFunc(out, threadsCount)
	// })
	time.Sleep(5 * time.Second)
	finishFunc()
	wg.Wait()
	time.Sleep(1 * time.Second)
}
