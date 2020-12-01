package notifier

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	maxIdleConnections = 5
	requestTimeout     = 60 * time.Second
)

type WebhookNotifier struct {
	url    string
	client *http.Client
}

func NewWebhookNotifier(url string) Notifier {
	return &WebhookNotifier{
		url: url,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: maxIdleConnections,
			},
			Timeout: requestTimeout,
		},
	}
}

func (n *WebhookNotifier) Notify(t UsdtTransaction) {
	body, _ := json.Marshal(t)
	req, err := http.NewRequest("POST", n.url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Webhook notifier error. %+v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := n.client.Do(req)
	if err != nil && response == nil {
		log.Printf("Error sending request to webhook URL. %+v\n", err)
		return
	}

	defer response.Body.Close()
	ioutil.ReadAll(response.Body)
}
