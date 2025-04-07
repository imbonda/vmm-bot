# vmm-bot
A volume market-making bot built in Go and containerized with Docker.<br/>
The bot continuously places pseudo-random orders within the bid/ask spread to simulate liquidity.


## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/imbonda/vmm-bot.git
cd vmm-bot
```

### 2. Configure Environment

Create a `.env` file from the provided example:
```bash
cp .env.example .env
```

Edit `.env` file with api-keys and secrets:
```bash
# Bybit.
BYBIT_API_KEY= # bybit exchange api-key...
BYBIT_API_SECRET= # bybit exchange api-secret...

# Biconomy.
BICONOMY_API_KEY= # biconomy exchange api-key...
BICONOMY_API_SECRET= # biconomy exchange api-secret...
```

### 3. Run with Docker Compose

```bash
docker compose build
docker compose up -d
```

To stop:
```bash
docker compose stop
```

Thats it!

---

## 🌐 Supported Exchanges

| Exchange Name     | Value for `EXCHANGE_NAME`	 | Docs / Info                                             |
|-------------------|----------------------------|---------------------------------------------------------|
| Bybit             | `bybit`                    | [docs](https://bybit-exchange.github.io/docs/v5/intro)  |
| Biconomy          | `biconomy`                 | [docs](https://github.com/BiconomyOfficial/apidocs)     |

## 🧾 Enviornment Variables

| Variable                            | Description                                 | Example            |
|-------------------------------------|---------------------------------------------|--------------------|
| SERVICE_NAME                        | The name of your service                    | `bybit-vmm-bot`    |
| EXCHANGE_NAME                       | The exchange to trade on                    | `bybit`            |
| BYBIT_API_KEY                       | Bybit API key                               | `...`              |
| BYBIT_API_SECRET                    | Bybit API secret                            | `...`              |
| BICONOMY_API_KEY                    | Biconomy API key                            | `...`              |
| BICONOMY_API_SECRET                 | Biconomy API secret                         | `...`              |
| INTERVAL_EXECUTION_DURATION         | Interval duration                           | `30s`              |
| NUM_OF_TRADE_ITERATIONS_IN_INTERVAL | Number of trades per interval               | `3`                |
| CANDLE_HEIGHT                       | Price restriction as % of last price        | `0.005`            |
| SPREAD_MARGIN_MIN                   | `price >= bid + spread * min_margin`        | `0.2`              |
| SPREAD_MARGIN_MAX                   | `price <= bid + spread * max_margin`        | `0.8`              |
| TRADE_AMOUNT_MIN                    | `amount >= min`                             | `100`              |
| TRADE_AMOUNT_MAX                    | `amount <= max`                             | `200`              |
| SYMBOL                              | Trading pair symbol                         | `BTCUSDT`          |
| ORACLE_SYMBOL                       | Trading pair symbol used as price oracle    | `BTCUSDT`          |
| PRICE_DECIMALS_PRECISION            | Price decimals may differ by exchange       | `5`                |
| AMOUNT_DECIMALS_PRECISION           | Amount decimals may differ by exchange      | `3`                |

### 🔀 Trading Pair Symbol Format

| Exchange                  | Format Style                                 | Example           |
|---------------------------|---------------------------------------------|--------------------|
| Bybit                     | UPPERCASE, no separator                     | `BTCUSDT`          |
| Biconomy                  | UPPERCASE with underscore                   | `BTC_USDT`         |

### 🔀 Amount Decimals

| Exchange                  | Decimals    |
|---------------------------|-------------|
| Bybit                     | 0           |
| Biconomy                  | 3           |
