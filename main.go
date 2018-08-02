package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type request struct {
	URL     string      `json:"url"`
	Method  string      `json:"method"`
	Headers http.Header `json:"headers"`
	Body    []byte      `json:"body"`
}

const (
	loopCount      = 1
	delayWorkCount = 50000000
)

func doWork() {

	done := make(chan int)

	for i := 0; i < loopCount; i++ {
		select {
		case <-done:
			return
		default:
		}
	}

}

func handler(logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		logClientIP(logger, r)

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

	})
}

func delayHandler(logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		logClientIP(logger, r)

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

	})
}

func logClientIP(logger *log.Logger, r *http.Request) {

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return
	}

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		return
	}

	logger.Printf("Request received from %s", clientIP)
}

func main() {

	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("Server is starting...")

	http.Handle("/", handler(logger))
	http.Handle("/delay", delayHandler(logger))
	http.ListenAndServe(":8080", nil)
}
