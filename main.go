package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

const (
	workerCount = 2
	loopCount   = 200000000
)

func doWork() {

	done := make(chan int)

	var wg sync.WaitGroup

	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < loopCount/workerCount; i++ {
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

func handler(logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		logClientIP(logger, r)

		doWork()

		rw.Header().Set("Content-Type", "text/plain")
		rw.Write([]byte("Hello from goldengoose!"))

	})
}

func delayHandler(logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		logClientIP(logger, r)

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-time.After(10 * time.Second):
					return
				// Cause select loop to spin
				default:
				}
			}
		}()
		wg.Wait()

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

	r := httprouter.New()
	r.Handler("GET", "/", handler(logger))
	r.Handler("GET", "/delay", delayHandler(logger))

	http.ListenAndServe(":8000", r)
}
