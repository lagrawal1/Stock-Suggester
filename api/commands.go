package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func handlerIndustries(w http.ResponseWriter, req *http.Request) {
	if cfg.db == nil {
		os.Exit(1)
	}

	data, err := cfg.db.DistinctIndustries(context.Background())

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

func handlerSectors(w http.ResponseWriter, req *http.Request) {
	if cfg.db == nil {
		os.Exit(1)
	}

	data, err := cfg.db.DistinctSectors((context.Background()))

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

func handlerHighDivByIndustry(w http.ResponseWriter, req *http.Request) {

	var id int
	query_parameters := req.URL.Query()

	id_str := query_parameters.Get("industry_id")

	if id_str == "" {
		w.Write([]byte("industry_id query parameter not set"))
		w.WriteHeader(404)
	}

	id, err := strconv.Atoi(id_str)

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}

	data, err := cfg.db.BestDividendStocksByIndustry(context.Background(), int64(id))

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}

	if len(data) == 0 {
		w.Write([]byte("no data found"))
		w.WriteHeader(404)
	}

	data_mar, err := json.Marshal(data)

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}

	w.Write(data_mar)
	w.WriteHeader(http.StatusOK)
}

func handlerHighDivBySector(w http.ResponseWriter, req *http.Request) {

	var id int
	query_parameters := req.URL.Query()

	id_str := query_parameters.Get("sector_id")

	if id_str == "" {
		w.Write([]byte("sector_id query parameter not set"))
		w.WriteHeader(404)
	}

	id, err := strconv.Atoi(id_str)

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}
	data, err := cfg.db.HighDividendBySector(context.Background(), int64(id))

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}

	if len(data) == 0 {
		w.Write([]byte("no data found"))
		w.WriteHeader(404)
	}

	data_mar, err := json.Marshal(data)

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}

	w.Write(data_mar)
	w.WriteHeader(http.StatusOK)

}

func handlerHighFCFBySector(w http.ResponseWriter, req *http.Request) {

	var id int
	query_parameters := req.URL.Query()

	id_str := query_parameters.Get("sector_id")

	if id_str == "" {
		w.Write([]byte("sector_id query parameter not set"))
		w.WriteHeader(404)
	}

	id, err := strconv.Atoi(id_str)

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}
	data, err := cfg.db.HighCashFlowBySector(context.Background(), int64(id))

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}

	if len(data) == 0 {
		w.Write([]byte("no data found"))
		w.WriteHeader(404)
	}

	data_mar, err := json.Marshal(data)

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}

	w.Write(data_mar)
	w.WriteHeader(http.StatusOK)
}

func handlerHighGrowthBySector(w http.ResponseWriter, req *http.Request) {

	var id int
	query_parameters := req.URL.Query()

	id_str := query_parameters.Get("sector_id")

	if id_str == "" {
		w.Write([]byte("sector_id query parameter not set"))
		w.WriteHeader(404)
	}

	id, err := strconv.Atoi(id_str)

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}
	data, err := cfg.db.EarningsQuartGrowthBySector(context.Background(), int64(id))

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}

	if len(data) == 0 {
		w.Write([]byte("no data found"))
		w.WriteHeader(404)
	}

	data_mar, err := json.Marshal(data)

	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		w.WriteHeader(404)
	}

	w.Write(data_mar)
	w.WriteHeader(http.StatusOK)
}
