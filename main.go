package main

import (
	"Stock-Suggester/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func SrvReady(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
}

func StartDB() {
	err := godotenv.Load("data.env")
	if err != nil {
		fmt.Println(err)
		return
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("sqlite", dbURL)

	if err != nil {
		fmt.Println(err)
		return
	}

	cfg.db = database.New(db)
	fmt.Println(cfg.db)
}

func (cfg *Config) handlerIndustriesAPI(w http.ResponseWriter, req *http.Request) {
	if cfg.db == nil {
		os.Exit(1)
	}

	data, err := cfg.db.DistinctIndustry(context.Background())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jsondata, err := json.Marshal(&data)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	w.Write(jsondata)
	w.WriteHeader(http.StatusOK)

}

func API() {

	cfg = Config{}

	StartDB()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", SrvReady)

	mux.HandleFunc("GET /stocks", cfg.handlerIndustriesAPI)

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

func main() {
	CLI()
}
