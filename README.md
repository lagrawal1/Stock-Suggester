# Stock-Suggester : Submission for the BootDev 2025 Hackathon

## Description
A REPL that gives data anaylsis on a variety of NASDAQ stocks through data from the yfinance library! Get valuable data to improve your stock portfolio!

## Installation
To install Stock-Suggester, clone this repositiory and simply run "go run . " from the main directory. 

``` git clone https://github.com/lagrawal1/Stock-Suggester.git ``` <br>
``` go run . ```

## How to Use
Here is a list of commands that you can use!
-  exit : Exits the REPL.
-  help : Help command gives description of all commands.
-  industries : Gives a list of industries.
-  sectors : Gives a list of sectors

For these commands an integer must be given based on the sectors and industries.
-  highdivbyind : Gives the top 5 highest dividend stocks in a given industry.
-  highdivbysec : Gives the top 5 highest dividend stocks in a given sector.
-  highfcf : Gives the top 5 highest free cash flow stocks in a given sector.
-  growth : Gives the top 5 highest growth stocks by sector.

``` cmd_name id_for_industry_or_sector``` <br> Ex: ```highdivbyind 15```



## Technology and Tools
The api directory shows the api that I made to serve stock data from a server from my computer. The data goes through a cloudflare tunnel and serves to the REPL. If the api is down (since I'm sleepin), data from the local sqlite database will be used. 

Langauges: Go and SQL.<br>
SQLC is used to generate Go from SQL.<br>  
Database: SQLite.<br>
Serving API: Cloudflare Tunnels.<br>


## Notes

NOTE: ALL DATA USED IS FROM yfinance. I OWN NO RIGHTS TO THIS DATA. CONSIDER THIS DATA AS TEST DATA. THIS MAY NOT BE COMPLETELY ACCURATE AS DATA CHANGES QUICKLY IN THE STOCK MARKET.

