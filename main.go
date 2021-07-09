package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
			defer wg.Done()
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

	// Setup signal handler for TERM signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for range c {
			time.Sleep(3 * time.Second)
			os.Exit(0)
		}
	}()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Println("Server is starting...")

	r := httprouter.New()
	r.Handler("GET", "/", handler(logger))
	r.Handler("GET", "/delay", delayHandler(logger))

	http.ListenAndServe(":8000", r)
}
