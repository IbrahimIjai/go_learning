package utils

import "time"

const EventBufferSize = 100

// ChainEvent is the normalized shape used by every listener.
// EVM logs and SVM log notifications look different at the wire level, but the
// rest of the program should not care where the event came from.
type ChainEvent struct {
	Chain       string
	Runtime     string
	EventType   string
	TxHash      string
	BlockNumber uint64
	Amount      string
	Time        time.Time
	Raw         string
}

type Runtime string

const (
	RuntimeEVM Runtime = "evm"
	RuntimeSVM Runtime = "svm"
)

type ChainConfig struct {
	Name            string
	Runtime         Runtime
	WSURL           string
	ContractAddress string
	ProgramID       string
	EventTypes      []string
}
