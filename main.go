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
	// Work well with a 2 second timeout
	loopCount = 120000000
	// Works well with a 1 second timeout
	// loopCount = 60000000
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

func healthHandler(logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		logClientIP(logger, r)

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
	r.Handler("GET", "/healthz", healthHandler(logger))
	r.Handler("GET", "/delay", delayHandler(logger))

	http.ListenAndServe(":8000", r)
}
