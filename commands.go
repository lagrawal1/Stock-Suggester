package main

import (
	"Stock-Suggester/internal/database"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func validSectors(input int) error {
	if input <= 1 || input > 12 {
		return fmt.Errorf("Not valid id")

	}
	return nil
}

func validIndustries(input int) error {
	if input <= 0 || input > 146 || input == 4 {
		return fmt.Errorf("Not valid id")
	}
	return nil
}

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

func handlerIndustriesAPI() ([]database.Industry, error) {

	res, err := http.Get("https://www.lokics.xyz/stocks/industries")
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var data_unmar []database.Industry

	json.Unmarshal(data, &data_unmar)

	return data_unmar, nil
}

func handlerIndustries() error {

	if cfg.db == nil {
		fmt.Print("Start Database")
		os.Exit(1)
	}

	var data []database.Industry

	if err := APIHealth(); err != nil {

		fmt.Println("Using local database.")

		data_loc, err := cfg.db.DistinctIndustries(context.Background())

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data = append(data, data_loc...)
	} else {
		data_loc, err := handlerIndustriesAPI()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data = append(data, data_loc...)

	}

	fmt.Println("List of Industries")

	for _, industry := range data {

		fmt.Println("- ", industry.ID, ":", industry.Industry.String)

	}
	return nil

}

func handlerSectorsAPI() ([]database.Sector, error) {

	res, err := http.Get("https://www.lokics.xyz/stocks/sectors")
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var data_unmar []database.Sector

	json.Unmarshal(data, &data_unmar)

	return data_unmar, nil
}

func handlerSectors() error {

	if cfg.db == nil {
		fmt.Print("Start Database")
		os.Exit(1)
	}

	var data []database.Sector

	if err := APIHealth(); err != nil {

		fmt.Println("Using local database.")

		data_loc, err := cfg.db.DistinctSectors(context.Background())

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data = append(data, data_loc...)
	} else {
		data_loc, err := handlerSectorsAPI()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data = append(data, data_loc...)

	}

	fmt.Println("List of Sectors")

	for _, sector := range data {

		fmt.Println("- ", sector.ID, ":", sector.SectorName.String)

	}
	return nil

}

func handlerHighDivByIndAPI(id int64) ([]database.BestDividendStocksByIndustryRow, error) {

	res, err := http.Get("https://www.lokics.xyz/stocks/divbyindustry?industry_id=" + strconv.Itoa(int(id)))

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var data_unmar []database.BestDividendStocksByIndustryRow

	json.Unmarshal(data, &data_unmar)

	return data_unmar, nil
}

func handlerHighDivByInd() error {

	if len(cfg.args) < 1 {
		return fmt.Errorf("not enough arguments")
	}

	var industry_id int

	industry_id, err := strconv.Atoi(cfg.args[0])

	if err != nil {
		return err
	}

	err = validIndustries(industry_id)
	if err != nil {
		return err
	}

	data := make([]database.BestDividendStocksByIndustryRow, 0)
	if err := APIHealth(); err != nil {

		data_loc, err := cfg.db.BestDividendStocksByIndustry(context.Background(), int64(industry_id))

		if err != nil {
			return err
		}
		data = append(data, data_loc...)

		if len(data) == 0 {
			return fmt.Errorf("no data found")
		}
	} else {
		data_loc, err := handlerHighDivByIndAPI(int64(industry_id))

		if err != nil {
			return err
		}
		data = append(data, data_loc...)

		if len(data) == 0 {
			return fmt.Errorf("no data found")
		}
	}

	fmt.Println("\nTop Dividend Stocks in", data[0].Industry.String)

	for _, stock := range data {
		fmt.Println("- ", stock.Symbol.String, stock.Displayname.String, ":", stock.MaxDiv.Float64)
	}

	return nil

}

func handlerHighDivBySecAPI(id int64) ([]database.HighDividendBySectorRow, error) {

	res, err := http.Get("https://www.lokics.xyz/stocks/divbysector?sector_id=" + strconv.Itoa(int(id)))

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var data_unmar []database.HighDividendBySectorRow

	json.Unmarshal(data, &data_unmar)

	return data_unmar, nil
}

func handlerHighDivBySec() error {

	if len(cfg.args) < 1 {
		return fmt.Errorf("not enough arguments")
	}

	var sector_id int

	sector_id, err := strconv.Atoi(cfg.args[0])

	if err != nil {
		return err
	}

	err = validSectors(sector_id)
	if err != nil {
		return err
	}

	var data []database.HighDividendBySectorRow

	if err := APIHealth(); err != nil {

		data_loc, err := cfg.db.HighDividendBySector(context.Background(), int64(sector_id))
		if err != nil {
			return err
		}
		data = append(data, data_loc...)
		if len(data) == 0 {
			return fmt.Errorf("no data found")
		}
	} else {
		data_loc, err := handlerHighDivBySecAPI(int64(sector_id))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		data = append(data, data_loc...)
		if len(data) == 0 {
			return fmt.Errorf("no data found")
		}

	}

	fmt.Println("\nTop Dividend Stocks in", data[0].Sector.String)

	for _, stock := range data {
		fmt.Println("- ", stock.Symbol.String, stock.Displayname.String, ":", stock.MaxDiv.Float64)
	}

	return nil

}

func handlerHighFCFBySecAPI(id int64) ([]database.HighCashFlowBySectorRow, error) {

	res, err := http.Get("https://www.lokics.xyz/stocks/fcfbysector?sector_id=" + strconv.Itoa(int(id)))

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var data_unmar []database.HighCashFlowBySectorRow

	json.Unmarshal(data, &data_unmar)

	return data_unmar, nil
}

func handlerHighFCF() error {
	if len(cfg.args) < 1 {
		return fmt.Errorf("not enough arguments")
	}

	var id int

	id, err := strconv.Atoi(cfg.args[0])

	if err != nil {
		return err
	}
	err = validSectors(id)
	if err != nil {
		return err
	}

	var data []database.HighCashFlowBySectorRow

	if err := APIHealth(); err != nil {

		data_loc, err := cfg.db.HighCashFlowBySector(context.Background(), int64(id))
		if err != nil {
			return err
		}
		data = append(data, data_loc...)
		if len(data) == 0 {
			return fmt.Errorf("no data found")
		}
	} else {
		data_loc, err := handlerHighFCFBySecAPI(int64(id))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data = append(data, data_loc...)
		if len(data) == 0 {
			return fmt.Errorf("no data found")
		}
	}

	fmt.Println("\nTop Free Cash Flow Stocks in", data[0].Sector.String)

	for _, stock := range data {
		fmt.Println("- ", stock.Symbol.String, stock.Displayname.String, ":", formatIntegers(stock.Maxfcf.Int64))
	}

	return nil
}

func handlerHighGrowthBySecAPI(id int64) ([]database.EarningsQuartGrowthBySectorRow, error) {

	res, err := http.Get("https://www.lokics.xyz/stocks/growthbysector?sector_id=" + strconv.Itoa(int(id)))

	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var data_unmar []database.EarningsQuartGrowthBySectorRow

	json.Unmarshal(data, &data_unmar)

	return data_unmar, nil
}

func handlerHighGrowth() error {

	if len(cfg.args) < 1 {
		return fmt.Errorf("not enough arguments")
	}

	id, err := strconv.Atoi(cfg.args[0])

	if err != nil {
		return err
	}

	err = validSectors(id)
	if err != nil {
		return err
	}

	var data []database.EarningsQuartGrowthBySectorRow

	if err := APIHealth(); err != nil {

		data_loc, err := cfg.db.EarningsQuartGrowthBySector(context.Background(), int64(id))
		if err != nil {
			return err
		}
		data = append(data, data_loc...)
		if len(data) == 0 {
			return fmt.Errorf("no data found")
		}
	} else {
		data_loc, err := handlerHighGrowthBySecAPI(int64(id))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data = append(data, data_loc...)
		if len(data) == 0 {
			return fmt.Errorf("no data found")
		}

	}

	fmt.Println("\nTop Earnings Quarterly Growth Stocks in", data[0].SectorName.String)

	for _, stock := range data {
		fmt.Println("- ", stock.Symbol.String, stock.Displayname.String, ":", stock.Maxfcf.Float64)
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
