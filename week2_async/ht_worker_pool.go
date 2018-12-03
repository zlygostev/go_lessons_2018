package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func doJob(log io.Writer, threadID int, in <-chan MsgType) {
	fmt.Fprintf(log, "[%d] Worker  is started...\n", threadID)
	count := 0
	defer fmt.Fprintf(log, "[%d][%d] finished...\n", threadID, count)
	for newMsg := range in {
		fmt.Fprintf(log, "[%d][%d] Msg [%s]\n", threadID, count, newMsg.text)
		time.Sleep(2000 * time.Millisecond)
		fmt.Fprintf(log, "[%d][%d] wait...\n", threadID, count)
		count++
	}
	// for {
	// 	select {
	// 	case msg := <-in:
	// 		fmt.Fprintf(log, "[%d][%d] Msg [%s]\n", threadID, count, msg.text)
	// 	// default:
	// 	// 	Sleep
	// 	count++
	// }
	fmt.Fprintf(log, "[%d] Worker  is finished on [%d]...\n", threadID, count)
}

func periodicJob(log io.Writer, threadID int, secs int) {
	fmt.Fprintf(log, "[%d] Periodic Job is Started...\n", threadID)
	ticker := time.NewTicker(time.Duration(secs) * time.Second)
	counter := 0
	for tickerTime := range ticker.C {
		fmt.Fprintf(log, "[%d] Step %d at %v...\n", threadID, counter, tickerTime)
		counter++
		if counter > 2 {
			ticker.Stop()
			break
		}
	}
	fmt.Fprintf(log, "[%d] Periodic Job has been finished...\n", threadID)
}

func poolEventsGenerator(log io.Writer, id int, timeout int, out chan<- MsgType, control <-chan struct{}) {
	defer close(out)
	fmt.Fprintf(log, "[%d] Generator is started...\n", id)
	count := 0

LOOP:
	for {
		timer := time.NewTimer(time.Duration(timeout) * time.Millisecond)
		msg := MsgType{
			tp:   POST,
			text: "Message " + strconv.Itoa(count),
		}
		select {
		case out <- msg:
			//free timer resource
			if !timer.Stop() {
				<-timer.C
			}

			fmt.Fprintf(log, "[%d] Write [%d]\n", id, count)
			time.Sleep(time.Duration(timeout) * time.Millisecond)
		case _ = <-control:
			fmt.Fprintf(log, "[%d] Stop signal on %d\n", id, count)
			//free timer resource
			if !timer.Stop() {
				<-timer.C
			}
			break LOOP
		case <-timer.C:
			fmt.Fprintf(log, "[%d] Timeout occured on %d\n", id, count)
			//Couldn't free resources till event occur
			// case <-time.After(time.Duration(timeout) * time.Millisecond):
			// 	fmt.Fprintf(log, "[%d] Timeout occured on %d\n", id, count)
		}
		count++
	}
	fmt.Fprintf(log, "[%d] Generator is finished...\n", id)

}

func afterFunc(log io.Writer, threadID int) {
	fmt.Fprintf(log, "[%d] I want to say good bye...\n", threadID)
}

const (
	//GET message type id
	GET = 0
	//POST message type id
	POST = 1
	//PUT message type id
	PUT = 2
	//DELETE message type id
	DELETE = 3
)

//MsgType message
type MsgType struct {
	tp   int
	text string
}

const threadsCount int = 10

func main() {
	out := os.Stdout
	fmt.Fprintf(out, "%v\n", threadsCount)
	msgChannel := make(chan MsgType, 1)
	controlChannel := make(chan struct{})
	go poolEventsGenerator(out, 0, 2000, msgChannel, controlChannel)
	go periodicJob(out, 1, 2)
	for thID := 2; thID < threadsCount; {
		go doJob(out, thID, msgChannel)
		thID++
	}
	afterFuncTimer := time.AfterFunc(1*time.Second, func() {
		afterFunc(out, threadsCount)
	})

	time.Sleep(5 * time.Second)
	controlChannel <- struct{}{}
	time.Sleep(1 * time.Second)
	afterFuncTimer.Stop()
}
