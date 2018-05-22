package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"
	"time"
)

type request struct {
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

const (
	loopCount      = 50000000
	delayWorkCount = 50000000
)

func doWork() {

	var wg sync.WaitGroup
	done := make(chan int)
	numCPU := runtime.NumCPU()

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < loopCount; j++ {
				select {
				case <-done:
					return
				default:
				}
			}

		}()
	}

	wg.Wait()
}

func handle(rw http.ResponseWriter, r *http.Request) {

	doWork()

	var err error
	rr := &request{}
	rr.Method = r.Method
	rr.Headers = r.Header
	rr.URL = r.URL.String()
	rr.Body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rrb, err := json.Marshal(rr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(rrb)
}

func delayHandle(rw http.ResponseWriter, r *http.Request) {

	done := make(chan int)

	go func() {
		for {
			select {
			case <-done:
			default:
			}
		}
	}()

	time.Sleep(10 * time.Second)
	close(done)

}

func main() {
	http.HandleFunc("/", handle)
	http.HandleFunc("/delay", delayHandle)
	http.ListenAndServe(":8000", nil)
}
