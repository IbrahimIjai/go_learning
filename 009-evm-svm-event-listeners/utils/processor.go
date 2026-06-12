package utils

import (
	"context"
	"fmt"
	"sync"
)

func ProcessEvents(ctx context.Context, events <-chan ChainEvent, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Event processor started")

	for {
		select {
		case <-ctx.Done():
			drainEvents(events)
			fmt.Println("Event processor stopping")
			return
		case event, ok := <-events:
			if !ok {
				fmt.Println("Event processor stopping")
				return
			}
			printEvent(event)
		}
	}
}

func drainEvents(events <-chan ChainEvent) {
	for {
		select {
		case event, ok := <-events:
			if !ok {
				return
			}
			printEvent(event)
		default:
			return
		}
	}
}

func printEvent(event ChainEvent) {
	fmt.Printf(
		"[%s] Runtime=%s | Chain=%s | Event=%s | Block/Slot=%d | Amount=%s | Tx=%s\n",
		event.Time.Format("15:04:05"),
		event.Runtime,
		event.Chain,
		event.EventType,
		event.BlockNumber,
		defaultString(event.Amount, "n/a"),
		event.TxHash,
	)
}
