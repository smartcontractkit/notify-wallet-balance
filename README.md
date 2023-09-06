# Notify Wallet Balance

[![Linting](https://github.com/smartcontractkit/notify-wallet-balance/actions/workflows/lint.yml/badge.svg)](https://github.com/smartcontractkit/notify-wallet-balance/actions/workflows/lint.yml)

A basic micro service to track wallet addresses on different chains, and notify based on low balance thresholds.

## Config

Configuration is done entirely through environment variables.

The app can notify you through slack. To do so, you need to install your own slack bot, and use its auth.

```bash
SLACK_API_KEY="xoxb-abc" ## Slack API key for bot
SLACK_CHANNEL="C123" ## Slack Channel to send messages to
SLACK_USER="U111" ## Slack User ID to notify for slack messages
NOTIFICATION_INTERVAL="4h" ## How often to check and notify
```

Set which networks you would like to check on with `NETWORK_PREFIXES`.

```bash
NETWORK_PREFIXES="OPTIMISM_GOERLI,ARBITRUM_GOERLI"
```

Each Network listed is used as a prefix for other env vars.

```bash
# Arbitrum Goerli
ARBITRUM_GOERLI_URL="wss://goerli-url" ## Websocket URL for the network
ARBITRUM_GOERLI_ADDRESSES="0xaaa,0xbbb" ## List of addresses to monitor
ARBITRUM_GOERLI_LOWER_LIMIT=30 ## How many ETH to consider worth notifying about a low balance
ARBITRUM_GOERLI_EXPLORER_URL="https://goerli.arbiscan.io/" ## URL to explorer if available for convenient links

# Optimism Goerli
OPTIMISM_GOERLI_URL="wss://optimism-url"
OPTIMISM_GOERLI_ADDRESSES="0xaaa,0xbbb"
OPTIMISM_GOERLI_LOWER_LIMIT=10
OPTIMISM_GOERLI_EXPLORER_URL="https://goerli-optimism.etherscan.io/"
```

## Run

When running locally, you can use plain go.

```sh
go run .
```

A docker image can be found at [kalverra/notify-wallet-balance](https://hub.docker.com/repository/docker/kalverra/notify-wallet-balance/general).
