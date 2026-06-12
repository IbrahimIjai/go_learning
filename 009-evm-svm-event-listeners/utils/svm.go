package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func listenToSVM(ctx context.Context, chain ChainConfig, events chan<- ChainEvent) {
	fmt.Printf("Started SVM listener for %s\n", chain.Name)

	runWithReconnect(ctx, chain.Name, func(attemptCtx context.Context) error {
		conn, _, err := websocket.DefaultDialer.DialContext(attemptCtx, chain.WSURL, http.Header{})
		if err != nil {
			return err
		}
		defer conn.Close()

		go func() {
			<-attemptCtx.Done()
			_ = conn.Close()
		}()

		if err := conn.WriteJSON(svmLogsSubscribeRequest(chain.ProgramID)); err != nil {
			return err
		}

		for {
			if err := attemptCtx.Err(); err != nil {
				return err
			}

			_, payload, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			event, ok := svmMessageToEvent(chain, payload)
			if !ok {
				continue
			}
			publishEvent(attemptCtx, events, event)
		}
	})
}

func svmLogsSubscribeRequest(programID string) map[string]any {
	return map[string]any{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "logsSubscribe",
		"params": []any{
			map[string]any{"mentions": []string{programID}},
			map[string]any{"commitment": "confirmed"},
		},
	}
}

func svmMessageToEvent(chain ChainConfig, payload []byte) (ChainEvent, bool) {
	var msg struct {
		Method string `json:"method"`
		Params struct {
			Result struct {
				Context struct {
					Slot uint64 `json:"slot"`
				} `json:"context"`
				Value struct {
					Signature string   `json:"signature"`
					Err       any      `json:"err"`
					Logs      []string `json:"logs"`
				} `json:"value"`
			} `json:"result"`
		} `json:"params"`
	}
	if err := json.Unmarshal(payload, &msg); err != nil || msg.Method != "logsNotification" {
		return ChainEvent{}, false
	}
	if msg.Params.Result.Value.Err != nil {
		return ChainEvent{}, false
	}

	rawLogs := strings.Join(msg.Params.Result.Value.Logs, " | ")
	eventType := inferSVMEventType(chain.EventTypes, rawLogs)

	return ChainEvent{
		Chain:       chain.Name,
		Runtime:     string(chain.Runtime),
		EventType:   eventType,
		TxHash:      msg.Params.Result.Value.Signature,
		BlockNumber: msg.Params.Result.Context.Slot,
		Amount:      "",
		Time:        time.Now(),
		Raw:         rawLogs,
	}, true
}

func inferSVMEventType(eventTypes []string, logs string) string {
	lowerLogs := strings.ToLower(logs)
	for _, eventType := range eventTypes {
		if strings.Contains(lowerLogs, strings.ToLower(eventType)) {
			return eventType
		}
	}
	return "ProgramLog"
}
