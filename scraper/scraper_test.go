package scraper_test

import (
	"testing"

	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/tpanum/metalligaen/scraper"
)

const (
	NEGOTIATE_URI = "/signalr/negotiate?clientProtocol=1.5&connectionData=%5B%7B%22name%22%3A%22sportsadminlivehub%22%7D%5D"
	CONNECT_URI   = "/signalr/connect?clientProtocol=1.5&connectionData=%5B%7B%22name%22%3A%22sportsadminlivehub%22%7D%5D&connectionToken=N69WJN8uyUHbcWfi8&tid=1&transport=serverSentEvents"
	CONNTOKEN     = "N69WJN8uyUHbcWfi8"
)

var (
	reqResp = map[string]string{
		NEGOTIATE_URI: fmt.Sprintf(`{"ConnectionToken":"%s"}`, CONNTOKEN),
	}
)

func MetalligaenStub() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if resp, ok := reqResp[r.RequestURI]; ok {
			w.Write([]byte(resp))
			return
		}

		if r.RequestURI == CONNECT_URI {
			w.Header().Add("Content-Type", "text/event-stream")
			flush, _ := w.(http.Flusher)
			flush.Flush()
			hj, _ := w.(http.Hijacker)
			conn, _, _ := hj.Hijack()
			for i := 0; ; i++ {
				if _, err := fmt.Fprintf(conn, "data: %d\n\n", i); err != nil {
					fmt.Println("write: ", err)
					conn.Close()
					return
				}
				fmt.Println(i)
				time.Sleep(1000 * time.Millisecond)
			}
		}

		http.Error(w, "not found", http.StatusNotFound)
	}))
}

func TestNewClient(t *testing.T) {
	server := MetalligaenStub()
	defer server.Close()

	if _, err := scraper.NewClient(server.URL); err != nil {
		t.Error("Unable to start scraper client. Error: " + err.Error())
		return
	}
}
