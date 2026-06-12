# Multichain Event Listener

A small Go project for learning goroutines, channels, WebSocket listeners, and graceful shutdown with real EVM and SVM chains.

## Mental Model

```text
Real EVM listener ┐
Real SVM listener ├──> events channel ──> processor
Real EVM listener ┘
```

Each listener runs in its own goroutine.

All listeners send events into one shared channel.

The processor reads from that channel in one place.

## Real EVM

```bash
export ETHEREUM_WSS_URL="wss://your-ethereum-rpc"
export ETHEREUM_CONTRACT_ADDRESS="0xYourContractAddress"

go run .
```

Optional Base listener:

```bash
export BASE_WSS_URL="wss://your-base-rpc"
export BASE_CONTRACT_ADDRESS="0xYourContractAddress"

go run .
```

## Real SVM / Solana

```bash
export SOLANA_WSS_URL="wss://api.mainnet-beta.solana.com"
export SOLANA_PROGRAM_ID="YourProgramPublicKey"

go run .
```

## Stop

```bash
Ctrl+C
```

## Project Layout

```text
main.go              app startup and shutdown
utils/config.go      chain configuration
utils/listeners.go   listener routing and reconnects
utils/evm.go         EVM log subscription
utils/svm.go         SVM logs subscription
utils/processor.go   shared event processor
utils/types.go       shared types
utils/helpers.go     small helpers
```
