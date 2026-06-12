package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ibrahimijai/multichain-event-listener/utils"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	events := make(chan utils.ChainEvent, utils.EventBufferSize)
	chains := utils.LoadChainConfigs()
	if len(chains) == 0 {
		fmt.Println("No real listeners configured.")
		fmt.Println("Set EVM or SVM environment variables, then run: go run .")
		return
	}

	var listeners sync.WaitGroup
	var processor sync.WaitGroup

	processor.Add(1)
	go utils.ProcessEvents(ctx, events, &processor)

	for _, chain := range chains {
		listeners.Add(1)
		go utils.RunListener(ctx, chain, events, &listeners)
	}

	fmt.Println("Multichain event listener is running")
	fmt.Println("Press Ctrl+C to stop")

	<-ctx.Done()
	fmt.Println("\nShutdown signal received...")

	listeners.Wait()
	close(events)
	processor.Wait()

	fmt.Println("Shutdown complete")
}
