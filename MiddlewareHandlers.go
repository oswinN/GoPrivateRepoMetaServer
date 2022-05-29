package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"runtime"
	"strings"
	"time"

	"github.com/throttled/throttled"
	"github.com/throttled/throttled/store/memstore"
)

type Handlers struct {
	debugmode bool
}

func MakeHandlers(dbg bool) (retval *Handlers) {
	retval = &Handlers{debugmode: dbg}
	return
}

func Handle404(w http.ResponseWriter, rq *http.Request) {
	http.Error(w, "{'Error':'404 You are lost'}", http.StatusNotFound)
}

// timeout requests
func (h Handlers) TimeoutHandler(next http.Handler) http.Handler {
	if h.debugmode {
		return http.TimeoutHandler(next, 600*time.Second, "API timed out")
	} else {
		return http.TimeoutHandler(next, 5*time.Second, "API timed out")
	}
}

// log requests with timestamps
func (h Handlers) LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var t1, t2 time.Time
		if h.debugmode {
			log.Println("############ SERVING REQUEST: " + r.Method + " " + r.URL.String())
			log.Println("FULL REQUEST: ")
			// Save a copy of this request for debugging.
			requestDump, err := httputil.DumpRequest(r, false)
			if err != nil {
				log.Println(err)
			}
			log.Println(string(requestDump))
			t1 = time.Now()
		}
		next.ServeHTTP(w, r)
		if h.debugmode {
			t2 = time.Now()
			log.Printf("############ DONE SERVING [%s] %q in %v\n", r.Method, r.URL.String(), t2.Sub(t1))
		}
	}
	return http.HandlerFunc(fn)
}

// make sure panics are recovered
func (h Handlers) RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				if h.debugmode {
					log.Printf("panic: %+v", err)
					log.Printf(CallerInfo(1))
					log.Printf(CallerInfo(2))
					log.Printf(CallerInfo(3))
					log.Printf(CallerInfo(4))
					log.Printf(CallerInfo(5))
					log.Printf(CallerInfo(6))
					log.Printf(CallerInfo(7))
					log.Printf(CallerInfo(8))
					log.Printf(CallerInfo(9))
				}
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetRateLimiterHandler() func(h http.Handler) http.Handler {
	// set up throttling
	store, err := memstore.New(65536)
	if err != nil {
		log.Fatal(err)
	}
	quota := throttled.RateQuota{throttled.PerMin(90), 180}
	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		log.Fatal(err)
	}
	httpRateLimiter := throttled.HTTPRateLimiter{
		RateLimiter: rateLimiter,
		VaryBy:      &throttled.VaryBy{Path: true},
	}
	return httpRateLimiter.RateLimit
}

func CallerInfo(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, file, line, ok := runtime.Caller(depth)
	if ok {
		return fmt.Sprintf("Callstack Info: File: %s  Function: %s Line: %d", GetFileNameFromPath(file), runtime.FuncForPC(function).Name(), line)
	}
	return fmt.Sprintf("CallerInfo(): Can't determine callstack.")

}

// return source filename after the last slash
func GetFileNameFromPath(path string) string {
	index := strings.LastIndex(path, "/")
	if index == -1 {
		return path
	}
	return path[index+1:]

}
