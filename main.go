package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"gopkg.in/yaml.v2"
)

type EventListener struct {
	wsClient *model.WebSocketClient
	config   EventConfig
}

func NewEventListener(config EventConfig) (*EventListener, error) {
	ls := &EventListener{
		config: config,
	}
	var err *model.AppError
	ls.wsClient, err = model.NewWebSocketClient4(config.MatterMostWebsocketURL, config.MatterMostToken)
	if err != nil {
		return nil, err
	}
	return ls, nil
}

// Listen listens and handles incoming messages until the server dies
func (ls *EventListener) Listen() error {
	err := ls.wsClient.Connect()

	if err != nil {
		return err
	}

	ls.wsClient.Listen()

	httpClient := http.DefaultClient
	ctx := context.Background()

	for evt := range ls.wsClient.EventChannel {
		fmt.Printf("Got event: %+v\n", evt)
		if stringSliceContains(ls.config.MatterMostEvents, evt.EventType()) {
			eventJSONBuf, err := json.Marshal(evt)
			if err != nil {
				log.Printf("Could not json serialize event: %v", err)
				continue
			}
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, ls.config.URL, strings.NewReader(string(eventJSONBuf)))
			if err != nil {
				log.Printf("Could not create request: %+v", err)
				continue
			}
			resp, err := httpClient.Do(req)
			if err != nil {
				log.Printf("Could not perform request: %+v", err)
				continue
			}
			if resp.StatusCode == http.StatusOK {
				log.Printf("Request successful")
			} else {
				bodyBuf, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Could request body: %+v", err)
					continue
				}
				log.Printf("Non HTTP OK response: %d: %s", resp.StatusCode, bodyBuf)
			}
		}
	}
	return nil
}

func main() {
	config := Config{}
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%+v\n", config)

	// Starting all listeners
	for botName, botConfig := range config {
		fmt.Printf("Preparing %s\n", botName)
		el, err := NewEventListener(botConfig)
		if err != nil {
			log.Fatalf("Failed start start event listener: %v", err)
		}
		go el.Listen()
	}

	// Hack: temporary way to keep this running while bot listeners run in go routines
	for {
		time.Sleep(10 * time.Second)
	}
}
