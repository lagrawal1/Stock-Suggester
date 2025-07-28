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
}

func main() {
	CommandMap = make(map[string]Command)
	cfg = Config{}
	StartDB()

	CommandRegister("industries", "Gives a list of industries.", handlerIndustries)
	CommandRegister("sectors", "Gives a list of sectors", handlerSectors)
	CommandRegister("highdivbyind", "Gives the top 5 highest dividend stocks in a given industry.", handlerHighDivByInd)
	CommandRegister("highdivbysec", "Gives the top 5 highest dividend stocks in a given sector.", handlerHighDivBySec)
	CommandRegister("highfcf", "Gives the top 5 highest free cash flow stocks in a given sector.", handlerHighFCF)
	CommandRegister("exit", "Exits the REPL.", handlerExit)
	CommandRegister("growth", "Gives the top 5 highest growth stocks by sector.", handlerHighGrowth)
	CommandRegister("help", "Help command gives description of all commands.", handlerHelp)

	for {

		buf := bufio.NewScanner(os.Stdin)
		fmt.Print("Stock Suggester > ")
		buf.Scan()

		tokens := strings.Fields(strings.ToLower(buf.Text()))

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
