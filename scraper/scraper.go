package scraper

import (
	"fmt"
	"sync"

	"encoding/json"
	"net/http"
	"net/url"

	esource "github.com/donovanhide/eventsource"
)

type event struct {
	Event string
	Data  json.RawMessage
}

type Client struct {
	config        ClientConfig
	quit          chan struct{}
	eventHandlers map[string]chan chan json.RawMessage
	mutex         sync.Mutex

	isClosed bool
}

const (
	CONNTOKEN = "ConnectionToken"
	DOMAIN    = "http://metalligaen.dk"
)

type ClientConfig struct {
	Domain string
	Token  string
}

func NewClientWithConfig(conf ClientConfig) (*Client, error) {
	c := &Client{
		config:        conf,
		eventHandlers: make(map[string]chan chan json.RawMessage),
		isClosed:      true,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

	return c, nil
}

func NewClient() (*Client, error) {
	c := &Client{
		config: ClientConfig{
			Domain: DOMAIN,
		},
		eventHandlers: make(map[string]chan chan json.RawMessage),
		isClosed:      true,
	}

	if err := c.connect(); err != nil {
		return nil, err
	}

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
	v := url.Values{
		"clientProtocol": {"1.5"},
		"connectionData": {"[{\"name\":\"sportsadminlivehub\"}]"},
	}

	resp, err := http.Get(c.config.Domain + "/signalr/negotiate?" + v.Encode())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	m := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return err
	}

	if _, ok := m[CONNTOKEN]; !ok {
		return fmt.Errorf("Unable to retrieve connection token")
	}

	token, ok := m[CONNTOKEN].(string)
	if !ok {
		return fmt.Errorf("Unable to retrieve connection token")
	}

	c.config.Token = token

	v = url.Values{
		"transport":       {"serverSentEvents"},
		"clientProtocol":  {"1.5"},
		"connectionData":  {"[{\"name\":\"sportsadminlivehub\"}]"},
		"connectionToken": {c.config.Token},
		"tid":             {"1"},
	}

	quit := make(chan struct{})

	stream, err := esource.Subscribe(c.config.Domain+"/signalr/connect?"+v.Encode(), "")
	if err != nil {
		return err
	}

	resp, err = http.Get(c.config.Domain + "/signalr/start?" + v.Encode())
	if err != nil {
		return err
	}
	resp.Body.Close()

	go func() {

		for {
			select {
			case <-quit:
				return
			case err := <-stream.Errors:
				fmt.Println("Error : ", err)
			case e := <-stream.Events:
				var resp angularStruct
				if err := json.Unmarshal([]byte(e.Data()), &resp); err != nil {
					continue
				}

				if len(resp.M) == 0 {
					continue
				}

				c.receiveEvent(&event{
					Event: resp.M[0].M,
					Data:  resp.M[0].A[0],
				})
			}

		}
	}()

	c.quit = quit

	return nil
}

func (c *Client) HookEvent(name string) chan json.RawMessage {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	out := make(chan json.RawMessage, 1)

	echan, ok := c.eventHandlers[name]
	if !ok {
		echan = make(chan chan json.RawMessage, 15)
		c.eventHandlers[name] = echan
	}

	echan <- out

	return out
}

func (c *Client) receiveEvent(e *event) {
	fmt.Println("Got Event: ", e.Event)

	if clients, ok := c.eventHandlers[e.Event]; ok {
		if len(clients) == 0 {
			return
		}

		for c := range clients {
			c <- e.Data
		}

		return
	}
}

func (c *Client) Close() {
	c.quit <- struct{}{}
}
