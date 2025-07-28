package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/time/rate"
	_ "modernc.org/sqlite"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(10, 50)
		visitors[ip] = limiter
	}
	return limiter
}

func getIP(r *http.Request) string {
	ipPort := r.RemoteAddr
	ip := strings.Split(ipPort, ":")[0]
	return ip
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
func SrvReady(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func main() {

	cfg = Config{}

	StartDB()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", SrvReady)
	mux.HandleFunc("GET /stocks/industries", handlerIndustries)
	mux.HandleFunc("GET /stocks/sectors", handlerSectors)
	mux.HandleFunc("GET /stocks/divbyindustry", handlerHighDivByIndustry)
	mux.HandleFunc("GET /stocks/divbysector", handlerHighDivBySector)
	mux.HandleFunc("GET /stocks/fcfbysector", handlerHighFCFBySector)
	mux.HandleFunc("GET /stocks/growthbysector", handlerHighGrowthBySector)

	handler := rateLimitMiddleware(mux)

	server := http.Server{
		Handler: handler,
		Addr:    ":8080",
	}

	err := http.ListenAndServe(server.Addr, rateLimitMiddleware(server.Handler))

	if err != nil {
		fmt.Println(err)
		return
	}

}
