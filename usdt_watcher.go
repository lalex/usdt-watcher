package usdt_watcher

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/lalex/usdt-watcher/notifier"
)

const (
	usdtContract  = "0xdac17f958d2ee523a2206206994597c13d831ec7"
	transferTopic = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
)

type UsdtWatcher struct {
	url          string
	watchAddress []common.Hash
	notifiers    []notifier.Notifier
}

func New(websocketUrl string) *UsdtWatcher {
	return &UsdtWatcher{
		url: websocketUrl,
	}
}

func (w *UsdtWatcher) AddAddress(address string) {
	w.watchAddress = append(w.watchAddress, common.HexToHash(address))
}

func (w *UsdtWatcher) RegisterNotifier(n notifier.Notifier) {
	w.notifiers = append(w.notifiers, n)
}

func (w *UsdtWatcher) Run(ctx context.Context) {
	if len(w.watchAddress) == 0 {
		log.Fatal("Usdt-watcher failed to run on empty address list.")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Connect the client.
	client, err := rpc.Dial(w.url)
	if err != nil {
		log.Fatalf("Failed to connect to node: %+v", err)
	}
	logsCh := make(chan types.Log)

	// Subscription
	go func() {
		for {
			w.subscribeLogs(ctx, client, logsCh)
			time.Sleep(1 * time.Second)
		}
	}()

	// Read logs
	for {
		select {
		case log := <-logsCh:
			trans := fromLog(log)
			for _, n := range w.notifiers {
				go func(n notifier.Notifier) {
					n.Notify(trans)
				}(n)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (w *UsdtWatcher) subscribeLogs(ctx context.Context, client *rpc.Client, logsCh chan types.Log) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Build topic filter
	topicFilter := struct {
		Address common.Address
		Topics  [][]common.Hash
	}{
		Address: common.HexToAddress(usdtContract),
		Topics:  [][]common.Hash{{common.HexToHash(transferTopic)}, nil, w.watchAddress},
	}

	// Subscribe to new logs.
	sub, err := client.EthSubscribe(ctx, logsCh, "logs", topicFilter)
	if err != nil {
		return
	}

	// Wait subscribe to exit
	<-sub.Err()
}

func fromLog(log types.Log) notifier.UsdtTransaction {
	return notifier.UsdtTransaction{
		BlockNumber:     log.BlockNumber,
		TransactionHash: log.TxHash.Hex(),
		From:            common.BytesToAddress(log.Topics[1].Bytes()).Hex(),
		To:              common.BytesToAddress(log.Topics[2].Bytes()).Hex(),
		Amount:          common.BytesToHash(log.Data).Big().Uint64(),
		Removed:         log.Removed,
	}
}
