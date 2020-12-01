# usdt-watcher
Watch for Tether (USDt) transaction in Ethereum blockchain using JSON-RPC API through WebSocket.

# Usage
```go
package main

import (
    "context"

    usdt_watcher "github.com/lalex/usdt-watcher"
    notifier "github.com/lalex/usdt-watcher/notifier"
)

func main() {
    w := usdt_watcher.New("wss://mainnet.infura.io/ws/v3/_projectId_")
    w.AddAddress("0x3f5CE5FBFe3E9af3971dD833D26bA9b5C936f0bE")
    w.RegisterNotifier(notifier.NewConsoleNotifier())
    //w.RegisterNotifier(notifier.NewWebhookNotifier("https://webhook.site/_uuid_"))
    w.Run(context.Background())
}
```

# License
[Apache 2.0 License](LICENSE)