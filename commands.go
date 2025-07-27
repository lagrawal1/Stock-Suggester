package main

import (
	"Stock-Suggester/internal/database"
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Command struct {
	name        string
	description string
	handler     func() error
}

type Config struct {
	args []string
	db   *database.Queries
}

var cfg Config

var CommandMap map[string]Command

func CommandRegister(name string, description string, handler func() error) {
	CommandMap[name] = Command{name: name, description: description, handler: handler}
}

func handlerIndustries() error {
	if cfg.db == nil {
		fmt.Print("Start Database")
		os.Exit(1)
	}

	data, err := cfg.db.DistinctIndustry(context.Background())

	var Symbol sql.NullString
	Symbol.Scan("AAPL")
	data_2, err := cfg.db.GetStockDataBySymbol(context.Background(), Symbol)
	fmt.Println(data_2)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("List of Industries")
	industries := make([]string, 10)

	for _, industry := range data {

		industry_string := strings.TrimSpace(industry.String)
		if industry_string == "" || industry_string == "None" {
			print("hmm")
			continue
		}

		industries = append(industries, industry.String)

	}

	slices.Sort(industries)

	for _, industry := range industries {

		fmt.Println("- ", industry)

	}

	fmt.Print(industries[0])
	return nil

}

func handlerExit() error {

	os.Exit(0)
	return nil

}

func CLI() {

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
