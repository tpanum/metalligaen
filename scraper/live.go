package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/tpanum/metalligaen/models"
)

type rawTick struct {
	ID            uint   `json:"gameID"`
	GameClock     string `json:"GameClock"`
	TimeRunning   bool   `json:"isTimeRunning"`
	HomeTeamGoals int    `json:"homeTeamGoals"`
	AwayTeamGoals int    `json:"awayTeamGoals"`
}

var (
	clockParser = regexp.MustCompile(`([0-9]){2}:([0-9]){2}`)
)

func updateTimeFromTick(t rawTick) {

	fmt.Println("Trying to update...", t)

	matches := clockParser.FindStringSubmatch(t.GameClock)

	var duration uint
	if len(matches) == 2 {
		mins, _ := strconv.Atoi(matches[0])
		secs, _ := strconv.Atoi(matches[1])

		duration = uint(mins)*60 + uint(secs)
	}

	if m, ok := models.IDToMatch[t.ID]; ok {
		m.UpdateDuration(duration)
	}
}

func (c *Client) ListenLive() error {
	v := url.Values{
		"transport":       {"serverSentEvents"},
		"clientProtocol":  {"1.5"},
		"connectionData":  {"[{\"name\":\"sportsadminlivehub\"}]"},
		"connectionToken": {c.token},
		"tid":             {"1"},
	}

	postUrl := c.domain + "/signalr/send?" + v.Encode()
	data := strings.NewReader(url.Values{
		"data": {fmt.Sprintf(`{"H":"sportsadminlivehub","M":"Register","A":[%v,true,true],"I":0}`, 1326)},
	}.Encode())

	req, err := http.NewRequest("POST", postUrl, data)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	tickChan := c.HookEvent("tick")

	for {
		select {
		case data := <-tickChan:
			var t rawTick
			if err := json.Unmarshal(data, &t); err != nil {
				fmt.Println("ERROR ", err)
				continue
			}
			updateTimeFromTick(t)

			tickChan = c.HookEvent("tick")
		}
	}
}
