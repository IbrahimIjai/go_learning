package utils

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

func RunListener(ctx context.Context, chain ChainConfig, events chan<- ChainEvent, wg *sync.WaitGroup) {
	defer wg.Done()

	switch chain.Runtime {
	case RuntimeEVM:
		listenToEVM(ctx, chain, events)
	case RuntimeSVM:
		listenToSVM(ctx, chain, events)
	}
}

func runWithReconnect(ctx context.Context, chainName string, run func(context.Context) error) {
	backoff := time.Second

	for {
		if err := ctx.Err(); err != nil {
			fmt.Printf("Stopping listener for %s\n", chainName)
			return
		}

		err := run(ctx)
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("Stopping listener for %s\n", chainName)
			return
		}
		if ctx.Err() != nil {
			fmt.Printf("Stopping listener for %s\n", chainName)
			return
		}

		fmt.Printf("%s listener disconnected: %v. Reconnecting in %s\n", chainName, err, backoff)

		select {
		case <-ctx.Done():
			fmt.Printf("Stopping listener for %s\n", chainName)
			return
		case <-time.After(backoff):
		}

		if backoff < 30*time.Second {
			backoff *= 2
		}
	}
}

func publishEvent(ctx context.Context, events chan<- ChainEvent, event ChainEvent) {
	select {
	case <-ctx.Done():
	case events <- event:
	}
}
