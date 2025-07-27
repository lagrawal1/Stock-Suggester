package main

import (
	"Stock-Suggester/internal/database"
	"bufio"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

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

func main() {

	CommandMap = make(map[string]Command)
	cfg = Config{}
	StartDB()

	CommandRegister("industries", "za industries", handlerIndustries)

	for {

		cfg := Config{}
		buf := bufio.NewScanner(os.Stdin)
		fmt.Print("Stock Suggester > ")
		buf.Scan()

		tokens := strings.Fields(buf.Text())

		cmd := tokens[0]
		args := tokens[1:]

		cfg.args = args
		command, ok := CommandMap[cmd]

		if !ok {
			fmt.Println("Not a valid command")
		} else {
			command.handler()
		}
	}

}
