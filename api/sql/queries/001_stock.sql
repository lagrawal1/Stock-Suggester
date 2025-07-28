-- name: GetStockDataBySymbol :one
SELECT * FROM Stock WHERE symbol = ?;

-- name: BestDividendStocksByIndustry :many
SELECT 
    s1.symbol, 
    s1.displayName, 
    s1.industry, 
    s1.dividendYield AS max_div
FROM 
    Stock s1
JOIN 
    Industry i ON LOWER(s1.industry) = LOWER(i.industry)
WHERE 
    i.id = ?
    AND s1.industry IS NOT NULL
    AND TRIM(s1.industry) != ''
    AND (
        SELECT COUNT(*) 
        FROM Stock s2 
        WHERE LOWER(s2.industry) = LOWER(s1.industry)
          AND s2.dividendYield > s1.dividendYield
    ) < 5
ORDER BY 
    s1.dividendYield DESC 
LIMIT 5;

-- name: HighDividendBySector :many
SELECT 
    s1.symbol, 
    s1.displayName, 
    s1.sector, 
    s1.dividendYield AS max_div
FROM 
    Stock s1
JOIN 
    Sector sec ON LOWER(s1.sector) = LOWER(sec.sector_name)
WHERE 
    sec.id = ?
    AND s1.dividendYield IS NOT NULL
    AND (
        SELECT COUNT(*)
        FROM Stock s2
        WHERE LOWER(s2.sector) = LOWER(s1.sector)
          AND s2.dividendYield > s1.dividendYield
    ) < 5
ORDER BY 
    s1.dividendYield DESC
LIMIT 5;



-- name: HighCashFlowBySector :many
SELECT 
    s1.symbol, 
    s1.displayName, 
    s1.sector, 
    s1.freeCashflow AS maxFCF
FROM 
    Stock s1
JOIN 
    Sector sec ON LOWER(s1.sector) = LOWER(sec.sector_name)
WHERE 
    sec.id = ?
    AND (
        SELECT COUNT(*) 
        FROM Stock s2 
        WHERE LOWER(s2.sector) = LOWER(s1.sector)
          AND s2.freeCashflow > s1.freeCashflow
    ) < 5
ORDER BY 
    s1.freeCashflow DESC 
LIMIT 5;


-- name: EarningsQuartGrowthBySector :many
SELECT 
    s1.symbol, 
    s1.displayName, 
    sec.sector_name, 
    s1.earningsQuarterlyGrowth AS maxFCF
FROM 
    Stock s1
JOIN 
    Sector sec ON LOWER(s1.sector) = LOWER(sec.sector_name)
WHERE 
    sec.id = ?
    AND (
        SELECT COUNT(*) 
        FROM Stock s2 
        WHERE LOWER(s2.sector) = LOWER(s1.sector)
          AND s2.earningsQuarterlyGrowth > s1.earningsQuarterlyGrowth
    ) < 5
ORDER BY 
    s1.earningsQuarterlyGrowth DESC 
LIMIT 5;

-- name: DistinctIndustries :many
SELECT * FROM Industry;

-- name: DistinctSectors :many
SELECT * FROM Sector;




