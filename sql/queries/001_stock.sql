-- name: GetStockDataBySymbol :one
SELECT * FROM Stock WHERE symbol = ?;

-- name: DistinctIndustry :many
SELECT DISTINCT industry FROM Stock WHERE industry != 'None' AND industry IS NOT NULL ;

/*
-- name: BestDividendStocks :many
WITH HighestDiv AS (
    SELECT symbol, industry, MAX(dividend_rate) AS max_div, ROW_NUMBER() OVER (PARTITION BY industry ORDER BY dividend_rate) AS rn
    FROM Stock 
    WHERE industry != 'None' 
    GROUP BY (symbol, industry) 
    ORDER BY industry, max_div DESC)
SELECT symbol, industry, max_div FROM HighestDiv WHERE rn <= $1;

-- name: BestDividendStocksByIndustry :many
WITH HighestDiv AS (
    SELECT symbol, industry, MAX(dividend_rate) AS max_div, ROW_NUMBER() OVER (PARTITION BY industry ORDER BY dividend_rate) AS rn
    FROM Stock 
    WHERE industry != 'None' AND industry != '' 
    GROUP BY (symbol, industry) 
    ORDER BY industry, max_div DESC)
SELECT symbol, industry, max_div FROM HighestDiv WHERE rn <= $1 AND industry = $2;
*/

