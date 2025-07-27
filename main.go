package main

import (
	"Stock-Suggester/internal/database"
	"bufio"
	"database/sql"
	"fmt"
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

func main() {

	CommandMap = make(map[string]Command)
	cfg = Config{}
	StartDB()

	CommandRegister("industries", "za industries", handlerIndustries)
	CommandRegister("sectors", "za sectors", handlerSectors)
	CommandRegister("highdiv", "Gives the top 5 highest dividend stocks in a given industry", handlerHighDiv)
	CommandRegister("highfcf", "Gives the top 5 highest free cash flow stocks in a given sectore", handlerHighFCF)
	CommandRegister("exit", "Exits", handlerExit)

	for {

		buf := bufio.NewScanner(os.Stdin)
		fmt.Print("Stock Suggester > ")
		buf.Scan()

		tokens := strings.Fields(buf.Text())

		cmd := tokens[0]
		cfg.args = make([]string, 0)
		cfg.args = append(cfg.args, tokens[1:]...)
		command, ok := CommandMap[cmd]

		if !ok {
			fmt.Println("Not a valid command")
		} else {
			err := command.handler()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
