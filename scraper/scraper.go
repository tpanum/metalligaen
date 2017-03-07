package scraper

import (
	"fmt"
	"sync"

	"encoding/json"
	"net/http"
	"net/url"

	eclient "github.com/neelance/eventsource/client"
)

const METALID int = 1254

type Client struct {
	token     string
	quit      chan struct{}
	dataQueue chan json.RawMessage
	mutex     sync.Mutex

	isClosed bool
}

const (
	CONNTOKEN = "ConnectionToken"
	DOMAIN    = "http://metalligaen.dk"
)

func NewClient() (*Client, error) {
	v := url.Values{
		"clientProtocol": {"1.5"},
		"connectionData": {"[{\"name\":\"sportsadminlivehub\"}]"},
	}

	resp, err := http.Get(DOMAIN + "/signalr/negotiate?" + v.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	m := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return nil, err
	}

	if _, ok := m[CONNTOKEN]; !ok {
		return nil, fmt.Errorf("Unable to retrieve connection token")
	}

	token, ok := m[CONNTOKEN].(string)
	if !ok {
		return nil, fmt.Errorf("Unable to retrieve connection token")
	}

	c := &Client{
		token:    token,
		isClosed: true,
	}

	c.connect()

	return c, nil
}

type angularStruct struct {
	M []struct {
		H string
		M string
		A []json.RawMessage
	}
}

func (c *Client) connect() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	v := url.Values{
		"transport":       {"serverSentEvents"},
		"clientProtocol":  {"1.5"},
		"connectionData":  {"[{\"name\":\"sportsadminlivehub\"}]"},
		"connectionToken": {c.token},
		"tid":             {"1"},
	}

	client, err := eclient.New(DOMAIN + "/signalr/connect?" + v.Encode())
	if err != nil {
		return err
	}

	quit := make(chan struct{})
	dataQueue := make(chan json.RawMessage, 10)

	go func() {
		for {
			select {
			case <-quit:
				return
			case event := <-client.Stream:
				var resp angularStruct
				json.Unmarshal(event.Data, &resp)

				if len(resp.M) > 0 {
					dataQueue <- resp.M[0].A[0]
				}
			}
		}
	}()

	c.quit = quit
	c.dataQueue = dataQueue

	return nil
}

func (c *Client) Close() {
	c.quit <- struct{}{}
}
