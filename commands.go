package main

import (
	"Stock-Suggester/internal/database"
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

	var industry sql.NullString
	industry.Scan("Leisure")
	data_2, err := cfg.db.BestDividendStocksByIndustry(context.Background(), database.BestDividendStocksByIndustryParams{
		Industry:   industry,
		Industry_2: industry,
	},
	)

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

func handlerHighDiv() error {

	if len(cfg.args) < 1 {
		return fmt.Errorf("not enough arguments")
	}

	var industry_name sql.NullString

	industry_name.Scan(cfg.args[0])

	data, err := cfg.db.BestDividendStocksByIndustry(context.Background(), database.BestDividendStocksByIndustryParams{
		Industry:   industry_name,
		Industry_2: industry_name,
	})

	if err != nil {
		return err
	}

	fmt.Println("Top Dividend Stocks in", industry_name.String)

	for _, stock := range data {
		fmt.Println(stock.Symbol.String, ":", stock.MaxDiv.Float64)
	}

	return nil

}

func handlerExit() error {

	os.Exit(0)
	return nil

}
