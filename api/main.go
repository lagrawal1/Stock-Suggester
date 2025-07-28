package main

import (
	"fmt"
	"net/http"

	_ "modernc.org/sqlite"
)

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

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := http.ListenAndServe(server.Addr, server.Handler)

	if err != nil {
		fmt.Println(err)
		return
	}

}
