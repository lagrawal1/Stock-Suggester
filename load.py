import yfinance as yf
import pandas as pd
import sqlite3
from datetime import datetime

def detect_sqlite_type(val):
    """Return appropriate SQLite type for a Python value."""
    if isinstance(val, bool):
        return "INTEGER"
    elif isinstance(val, int):
        return "INTEGER"
    elif isinstance(val, float):
        return "REAL"
    elif isinstance(val, datetime):
        return "TEXT"
    elif isinstance(val, str):
        return "TEXT"
    elif val is None:
        return "TEXT"
    else:
        return "TEXT" 

def clean_value(val):
    if isinstance(val, (dict, list, set)):
        return str(val)
    elif isinstance(val, bool):
        return int(val)
    elif isinstance(val, datetime):
        return val.strftime('%Y-%m-%d %H:%M:%S')
    elif isinstance(val, (float, int, str)):
        return val
    elif val is None:
        return None
    else:
        return str(val)

def get_existing_columns(cursor, table):
    cursor.execute(f"PRAGMA table_info({table});")
    return {row[1] for row in cursor.fetchall()}

def create_table_if_not_exists(cursor, table):
    cursor.execute(f"""CREATE TABLE IF NOT EXISTS {table} (
        id INTEGER PRIMARY KEY AUTOINCREMENT
    );""")

def add_column(cursor, table, column, sqlite_type):
    try:
        cursor.execute(f'ALTER TABLE {table} ADD COLUMN "{column}" {sqlite_type};')
    except sqlite3.OperationalError as e:
        if "duplicate column name" in str(e):
            pass  # Already exists
        else:
            raise

def insert_dynamic(cursor, table, data_dict):
    existing_columns = get_existing_columns(cursor, table)

    # Add missing columns
    for key, val in data_dict.items():
        if key not in existing_columns:
            col_type = detect_sqlite_type(val)
            add_column(cursor, table, key, col_type)

    # Build insert
    columns = ", ".join([f'"{col}"' for col in data_dict.keys()])
    placeholders = ", ".join(["?"] * len(data_dict))
    values = [clean_value(val) for val in data_dict.values()]

    cursor.execute(
        f"INSERT INTO {table} ({columns}) VALUES ({placeholders});",
        values
    )

def main():
    column_names = ["Symbol", "Name", "Last Sale", "Net Change", "% Change", "Market Cap", "Country", "IPO Year", "Volume", "Sector", "Industry"]
    df = pd.read_csv("nasdaq_symbols.csv", usecols=column_names)
    symbols = df['Symbol']

    conn = sqlite3.connect("Stock_Suggester_DB")
    curr = conn.cursor()

    create_table_if_not_exists(curr, "Stock")

    for symbol in symbols:
        print(f"Processing: {symbol}")
        try:
            stock = yf.Ticker(symbol)
            info = stock.info

            if not info or "industry" not in info:
                print(f"No valid data for {symbol}")
                continue

            info_data = {}
            for k, v in info.items():
                # Optional: skip huge blobs or nested structures
                if isinstance(v, (dict, list, set)):
                    continue
                info_data[k] = v

            info_data["symbol"] = symbol  # Ensure symbol is always stored

            insert_dynamic(curr, "Stock", info_data)

        except Exception as e:
            print(f"Error processing {symbol}: {e}")
            continue

        conn.commit()

    curr.close()
    conn.close()

if __name__ == "__main__":
    main()
