-- name: GetStockDataBySymbol :one
SELECT * FROM Stock WHERE symbol = ?;

-- name: DistinctIndustry :many
SELECT DISTINCT industry FROM Stock WHERE industry != 'None' AND industry IS NOT NULL ;

-- name: DistinctSectors :many
SELECT DISTINCT sector FROM Stock WHERE sector != 'None' AND sector IS NOT NULL ;

-- name: BestDividendStocksByIndustry :many
SELECT s1.symbol, s1.industry, s1.dividendRate AS max_div
FROM Stock s1
WHERE s1.industry != 'None' 
  AND s1.industry != '' 
  AND s1.industry = ?
  AND (
      SELECT COUNT(*) 
      FROM Stock s2 
      WHERE s2.industry = s1.industry 
        AND s2.dividendRate > s1.dividendRate
  ) < ?
ORDER BY s1.dividendRate DESC LIMIT 5;

-- name: HighCashFlowBySector :many
SELECT s1.symbol, s1.sector, s1.freeCashflow AS maxFCF
FROM Stock s1
WHERE s1.sector != 'None' 
  AND s1.sector != '' 
  AND s1.sector = ?
  AND (
      SELECT COUNT(*) 
      FROM Stock s2 
      WHERE s2.sector = s1.sector 
        AND s2.freeCashflow > s1.freeCashflow
  ) < ?
ORDER BY s1.freeCashflow DESC LIMIT 5;



