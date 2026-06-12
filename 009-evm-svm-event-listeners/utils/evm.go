package utils

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func listenToEVM(ctx context.Context, chain ChainConfig, events chan<- ChainEvent) {
	fmt.Printf("Started EVM listener for %s\n", chain.Name)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(chain.ContractAddress)},
		Topics: [][]common.Hash{
			{crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))},
		},
	}

	runWithReconnect(ctx, chain.Name, func(attemptCtx context.Context) error {
		client, err := ethclient.DialContext(attemptCtx, chain.WSURL)
		if err != nil {
			return err
		}
		defer client.Close()

		logs := make(chan types.Log, 128)
		sub, err := client.SubscribeFilterLogs(attemptCtx, query, logs)
		if err != nil {
			return err
		}
		defer sub.Unsubscribe()

		for {
			select {
			case <-attemptCtx.Done():
				return attemptCtx.Err()
			case err := <-sub.Err():
				if err == nil {
					return errors.New("evm subscription closed")
				}
				return err
			case log := <-logs:
				publishEvent(attemptCtx, events, evmLogToEvent(chain, log))
			}
		}
	})
}

func evmLogToEvent(chain ChainConfig, log types.Log) ChainEvent {
	amount := ""
	if len(log.Data) >= 32 {
		amount = new(big.Int).SetBytes(log.Data[len(log.Data)-32:]).String()
	}

	raw := fmt.Sprintf("address=%s topics=%d data=0x%s", log.Address.Hex(), len(log.Topics), hex.EncodeToString(log.Data))

	return ChainEvent{
		Chain:       chain.Name,
		Runtime:     string(chain.Runtime),
		EventType:   "Transfer",
		TxHash:      log.TxHash.Hex(),
		BlockNumber: log.BlockNumber,
		Amount:      amount,
		Time:        time.Now(),
		Raw:         raw,
	}
}
