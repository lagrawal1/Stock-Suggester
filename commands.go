package main

import (
	"Stock-Suggester/internal/database"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
)

func APIHealth() error {
	res, err := http.Get("https://www.lokics.xyz/healthz")

	if err != nil {
		fmt.Println(err)
		return err
	}

	if res.StatusCode == 200 {
		return nil
	} else {
		if res.StatusCode == 530 {
			fmt.Println("API Server is currently offline.")
		}
		return fmt.Errorf("error: %v ", res.StatusCode)
	}
}

func formatIntegers(input int64) string {

	if input >= 1_000_000_000_000 {
		return fmt.Sprintf("%.1fT", float64(input)/1_000_000_000_000)
	} else if input >= 1_000_000_000 {
		return fmt.Sprintf("%.1fB", float64(input)/1_000_000_000)
	} else if input >= 1_000_000 {
		return fmt.Sprintf("%.1fM", float64(input)/1_000_000)
	}
	return fmt.Sprintf("%d", input)

}

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

func handlerIndustriesAPI() error {

	if err := APIHealth(); err != nil {
		return err
	}

	res, err := http.Get("https://www.lokics.xyz/stocks/industries")
	if err != nil {
		return err
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	var data_unmar []sql.NullString

	json.Unmarshal(data, &data_unmar)

	for _, industry := range data_unmar {
		fmt.Println(industry.String)
	}
	return nil
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
		fmt.Println(stock.Symbol.String, stock.Displayname.String, ":", stock.MaxDiv.Float64)
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

	fmt.Println("Top Free Cash Flow Stocks in", sector_name.String)

	for _, stock := range data {
		fmt.Println(stock.Symbol.String, stock.Displayname.String, ":", formatIntegers(stock.Maxfcf.Int64))
	}

	return nil

}

func handlerHighGrowth() error {

	if len(cfg.args) < 1 {
		return fmt.Errorf("not enough arguments")
	}

	var sector_name sql.NullString

	sector_name.Scan(cfg.args[0])

	data, err := cfg.db.EarningsQuartGrowthBySector(context.Background(), database.EarningsQuartGrowthBySectorParams{
		Sector:   sector_name,
		Sector_2: sector_name,
	})

	if err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("no data found")

	}

	fmt.Println("Top Earnings Quarterly Growth Stocks in", sector_name.String)

	for _, stock := range data {
		fmt.Println(stock.Symbol.String, stock.Displayname.String, ":", stock.Maxfcf.Float64)
	}

	return nil

}

func handlerExit() error {

	os.Exit(0)
	return nil
}

func handlerHelp() error {
	fmt.Println("Commands")

	for name, cmd := range CommandMap {

		fmt.Println("- ", name, ":", cmd.description)
	}
	return nil

}
