package scraper_test

import (
	"testing"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/donovanhide/eventsource"
	"github.com/tpanum/metalligaen/scraper"
)

type event struct {
	value string
}

func (e event) Id() string    { return "" }
func (e event) Event() string { return "message" }
func (e event) Data() string  { return e.value }

func returnStringValue(s string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(s))
	}
}

func pushValueToConn(r *http.Request, value string) {
	connToken := r.URL.Query().Get("connectionToken")
	if connToken == "" {
		return
	}

	eventServer.Publish([]string{connToken}, event{value: value})
}

func fileToString(file string) string {
	buf, _ := ioutil.ReadFile(file)
	return string(buf)
}

type angularRequest struct {
	M string
}

var (
	reqToData = map[string]string{
		"RegisterSchedule": fileToString("../test/register_schedule.json"),
	}
)

func sendToStream(w http.ResponseWriter, r *http.Request) {
	fmt.Println("0")

	data, ok := r.Form["data"]
	if !ok {
		return
	}

	fmt.Println("1")

	var req angularRequest
	if err := json.Unmarshal([]byte(data[0]), &req); err != nil {
		return
	}

	fmt.Println("2")

	if value, ok := reqToData[req.M]; ok {
		pushValueToConn(r, value)
	}

	w.Write([]byte(`{"I":"0"}`))
}

var (
	eventServer = eventsource.NewServer()
	conns       = make(map[string]bool)
)

func eventStream(w http.ResponseWriter, r *http.Request) {
	connToken := r.URL.Query().Get("connectionToken")
	if connToken == "" {
		return
	}

	conns[connToken] = true

	eventServer.Handler(connToken)(w, r)
	return
}

var reqResp = map[string]http.HandlerFunc{
	"/signalr/negotiate": returnStringValue(`{"ConnectionToken":"N69WJN8uyUHbcWfi8"}`),
	"/signalr/connect":   eventStream,
	"/signalr/start":     returnStringValue(`{"Response":"started"}`),
	"/signalr/send":      sendToStream,
}

func MetalligaenStub() *httptest.Server {

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI

		fmt.Println(r.RequestURI)

		for key, handler := range reqResp {
			if strings.HasPrefix(uri, key) {
				handler(w, r)
				return
			}
		}

		http.Error(w, "not found", http.StatusNotFound)
	}))
}

func TestNewClient(t *testing.T) {
	server := MetalligaenStub()
	defer server.Close()

	if _, err := scraper.NewClientWithConfig(scraper.ClientConfig{
		Domain: server.URL,
	}); err != nil {
		t.Error("Unable to start scraper client. Error: " + err.Error())
		return
	}
}

func TestRegisterSchedule(t *testing.T) {
	server := MetalligaenStub()
	defer server.Close()

	c, _ := scraper.NewClientWithConfig(scraper.ClientConfig{
		Domain: server.URL,
	})

	if _, err := c.GetSchedule(0); err != nil {
		t.Error("Unable to get schedule: " + err.Error())
		return
	}
}
