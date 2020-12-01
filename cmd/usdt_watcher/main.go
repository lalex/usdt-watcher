package main

import (
	"context"
	"os"

	usdt_watcher "github.com/lalex/usdt-watcher"
	notifier "github.com/lalex/usdt-watcher/notifier"
)

// Usage: ./usdt_watcher 0x5041ed759Dd4aFc3a72b8192C143F72f4724081A 0x3f5CE5FBFe3E9af3971dD833D26bA9b5C936f0bE

const (
	wsUrl = "wss://mainnet.infura.io/ws/v3/_projectid_"
)

func main() {
	w := usdt_watcher.New(wsUrl)
	for _, addr := range os.Args[1:] {
		w.AddAddress(addr)
	}
	w.RegisterNotifier(notifier.NewConsoleNotifier())
	//w.RegisterNotifier(notifier.NewWebhookNotifier("https://webhook.site/_uuid_"))
	w.Run(context.Background())
}
