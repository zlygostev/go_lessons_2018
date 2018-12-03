package main

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

// сюда писать код

type hash func(data string) string

var mu = &sync.Mutex{}

// GetMD5Optimally Thread safe function realization
func GetMD5Optimally(data string) string {
	mu.Lock()
	defer mu.Unlock()
	return DataSignerMd5(data)
}

func crcFunc(data string, out chan<- string) {
	crccode := DataSignerCrc32(data)
	out <- crccode
}

// func start(in, out chan interface{}) {

// 	for number := range in {
// 		var intNum int
// 		var ok bool
// 		if intNum, ok = number.(int); !ok {
// 			fmt.Printf("Income param is not int. %#v/n", number)
// 			panic("Income param is not int. %v/n")
// 		}

// 		strNum := strconv.Itoa(intNum)
// 		out <- strNum
// 	}
// }
type safeChan struct {
	waitedID   int
	maxInBufID int
	out        chan interface{}
	mtx        *sync.Mutex
	bufferData map[int]string
}

func createSafeWriteChan(out chan interface{}) safeChan {
	return safeChan{
		out:        out,
		mtx:        &sync.Mutex{},
		waitedID:   0,
		maxInBufID: 0,
		bufferData: make(map[int]string),
	}
}

func (ch *safeChan) write(val string, threadID int) {
	ch.mtx.Lock()
	defer ch.mtx.Unlock()
	if threadID == ch.waitedID {
		ch.out <- val
		ch.waitedID++
		ch.maxInBufID++
	} else {
		ch.bufferData[threadID] = val
		//	fmt.Println(threadID, "Add data in buffer", ". Wait", ch.waitedID)
		if threadID > ch.maxInBufID {
			ch.maxInBufID = threadID
			//		fmt.Println(threadID, "MaxID in buffer ", ch.maxInBufID, "Wait ", ch.waitedID)
		}
		return
	}
	//	fmt.Println(threadID, "Look through buffer. Max ", ch.maxInBufID, "Wait ", ch.waitedID)
	for val, ok := ch.bufferData[ch.waitedID]; ok && ch.waitedID < ch.maxInBufID; ch.waitedID++ {
		// if val, ok := ch.bufferData[ch.waitedID]; !ok {
		// 	//do something here
		// 	break
		// }
		fmt.Println(threadID, "Found postponed data for thread", ch.waitedID)
		ch.out <- val
		delete(ch.bufferData, ch.waitedID)
	}
}

func optmalSingleHashCalc(data string, threadID int, out *safeChan, wg *sync.WaitGroup) {
	defer wg.Done()
	t1 := time.Now()
	chan1 := make(chan string, 0)
	// Md5 + Crc32 channel implementation and run
	go func(data string, out chan<- string) {
		md5OfData := GetMD5Optimally(data)
		crcFunc(md5OfData, out)
	}(data, chan1)
	chan2 := make(chan string, 0)
	go crcFunc(data, chan2)
	// wait results
	var crcPart, md5CrcPart string
	var sumResult string
	for {
		select {
		case tmp := <-chan1:
			md5CrcPart = tmp
		case tmp := <-chan2:
			crcPart = tmp
			//TODO: timeout
		}
		if crcPart != "" && md5CrcPart != "" {
			sumResult = crcPart + "~" + md5CrcPart
			break
		}
	}
	out.write(sumResult, threadID)
	t2 := time.Now()
	fmt.Println(threadID, "Duration of one OptimalSingleHash calculation", t2.Sub(t1))
}

// SingleHash - parallel coding
func SingleHash(in, out chan interface{}) {
	// For each income data
	t00 := time.Now()
	threadID := 0
	thSafeChan := createSafeWriteChan(out)
	wg := &sync.WaitGroup{}
	for number := range in {
		var intNum int
		var ok bool
		if intNum, ok = number.(int); !ok {
			//fmt.Printf("Income param is not int. %#v/n", number)
			panic("Income param is not int. %#v/n")
		}

		strNum := strconv.Itoa(intNum)
		t1 := time.Now()
		wg.Add(1)
		go optmalSingleHashCalc(strNum, threadID, &thSafeChan, wg)
		threadID++
		t2 := time.Now()
		fmt.Println(threadID, "Duration of one SingleHash calculation", t2.Sub(t1))
	}
	//close(out)
	wg.Wait()
	t01 := time.Now()
	fmt.Println("Duration of all SingleHash calculation", t01.Sub(t00))
}
func threadHashCalculation(data string, threadID int, result *string, wg *sync.WaitGroup) {
	defer wg.Done()
	inData := strconv.Itoa(threadID) + data
	*result = DataSignerCrc32(inData)
}

func optmalMultiHashCalc(data string, threadID int, outCh *safeChan, wg *sync.WaitGroup) {
	defer wg.Done()
	wgInner := &sync.WaitGroup{}
	var results [6]string
	for i := 0; i < 6; i++ {
		wgInner.Add(1)
		go threadHashCalculation(data, i, &results[i], wgInner)
	}
	wgInner.Wait()
	out := ""
	for i := 0; i < 6; i++ {
		// if i > 0 {
		// 	out += "~"
		// }
		out += results[i]
	}
	//fmt.Println(threadID, "out: ", out)
	outCh.write(out, threadID)
}

//MultiHash function make multithiding calculations of crc32
func MultiHash(in, out chan interface{}) {
	// For each income data

	t0 := time.Now()
	threadID := 0
	thSafeChan := createSafeWriteChan(out)
	wg := &sync.WaitGroup{}

	for data := range in {
		strData, ok := data.(string)
		if !ok {
			fmt.Printf("Income param is not string. %v/n", data)
			panic("Income param is not int. %v/n")
		}
		fmt.Println(threadID, "Come ", strData)
		wg.Add(1)
		go optmalMultiHashCalc(strData, threadID, &thSafeChan, wg)
		threadID++
	}
	wg.Wait()
	t01 := time.Now()
	fmt.Println("Duration of all MultiHash calculation", t01.Sub(t0))
	//close(out)
}

//CombineResults function reduce results of worker and serialize it in one string
func CombineResults(in, out chan interface{}) {
	var array []string
	for data := range in {
		strData, ok := data.(string)
		if !ok {
			//fmt.Printf("Income param is not string. %v/n", data)
			panic("Income param is not int. %v/n")
		}
		array = append(array, strData)
	}
	sort.Strings(array)
	result := ""
	for i, str := range array {
		if i > 0 {
			result += "_"
		}
		result += str
	}
	out <- result
	//close(out)
}

//ExecutePipeline  conveer of income functions
func ExecutePipeline(funcs ...job) {

	in := make(chan interface{}, 0)
	out := make(chan interface{}, 0)
	wg := &sync.WaitGroup{}
	for _, jb := range funcs {
		in = out
		out = make(chan interface{}, 0)
		wg.Add(1)
		go func(in, out chan interface{}, wg *sync.WaitGroup, jb job) {
			defer wg.Done()
			defer close(out)
			jb(in, out)
		}(in, out, wg, jb)
	}
	wg.Wait()
	// //Read results of the last channel till the end
	// for result := range out {
	// 	fmt.Printf("result %s\n", result)
	// }

}
func main() {
	jobs := []job{}
	ExecutePipeline(jobs...)
}
