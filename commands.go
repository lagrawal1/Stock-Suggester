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

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("List of Industries")
	industries := make([]string, 0)

	for _, industry := range data {

		industries = append(industries, industry.String)

	}

	slices.Sort(industries)

	for _, industry := range industries {

		fmt.Println("- ", industry)

	}

	return nil

}

func handlerSectors() error {

	if cfg.db == nil {
		fmt.Print("Start Database")
		os.Exit(1)
	}

	data, err := cfg.db.DistinctSectors(context.Background())

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("List of Industries")
	sectors := make([]string, 0)

	for _, sector := range data {

		sector_string := strings.TrimSpace(sector.String)
		if sector_string == "" || sector_string == "None" {
			continue
		}

		sectors = append(sectors, sector.String)

	}

	slices.Sort(sectors)

	for _, sector := range sectors {

		fmt.Println("- ", sector)

	}
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

	if len(data) == 0 {
		return fmt.Errorf("no data found")

	}

	fmt.Println("Top Dividend Stocks in", industry_name.String)

	for _, stock := range data {
		fmt.Println(stock.Symbol.String, ":", stock.MaxDiv.Float64)
	}

	return nil

}

func handlerHighFCF() error {
	if len(cfg.args) < 1 {
		return fmt.Errorf("not enough arguments")
	}

	var sector_name sql.NullString

	sector_name.Scan(cfg.args[0])

	data, err := cfg.db.HighCashFlowBySector(context.Background(), database.HighCashFlowBySectorParams{
		Sector:   sector_name,
		Sector_2: sector_name,
	})

	if err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("no data found")

	}

	fmt.Println("Top Dividend Stocks in", sector_name.String)

	for _, stock := range data {
		fmt.Println(stock.Symbol.String, ":", stock.Maxfcf.Int64)
	}

	return nil

}

func handlerExit() error {

	os.Exit(0)
	return nil

}
