package scraper

import (
	"fmt"
	"strings"
	"time"

	"encoding/json"
	"net/http"
	"net/url"

	"github.com/tpanum/metalligaen/api"
)

func (c *Client) GetSchedule(id int) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	v := url.Values{
		"transport":       {"serverSentEvents"},
		"clientProtocol":  {"1.5"},
		"connectionData":  {"[{\"name\":\"sportsadminlivehub\"}]"},
		"connectionToken": {c.token},
	}

	postUrl := DOMAIN + "/signalr/send?" + v.Encode()
	data := strings.NewReader(url.Values{
		"data": {fmt.Sprintf(`{"H":"sportsadminlivehub","M":"RegisterSchedule","A":[%v],"I":0}`, id)},
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

	rawData := <-c.dataQueue
	var jsonData struct {
		Games []rawMatch `json:"schedule"`
	}

	if err := json.Unmarshal(rawData, &jsonData); err != nil {
		return err
	}

	matches := make([]*api.Match, len(jsonData.Games))
	for i, m := range jsonData.Games {
		timeOfMatch, err := time.Parse("2006-01-02T15:04", m.Date[0:11]+m.Time)
		if err != nil {
			fmt.Println(err)
		}

		hometeam, err := api.TeamTagFromName(m.HomeTeam)
		if err != nil {
			fmt.Println(err)
		}

		awayteam, err := api.TeamTagFromName(m.AwayTeam)
		if err != nil {
			fmt.Println(err)
		}

		matches[i] = &api.Match{
			ID:          uint(m.ID),
			HomeTeamTag: hometeam,
			AwayTeamTag: awayteam,
			TimeOfMatch: timeOfMatch,
		}
	}

	api.LoadMatches(matches)

	return nil
}

type rawMatch struct {
	ID           int    `json:"gameID"`
	Date         string `json:"gameDate"`
	Time         string `json:"gameTime"`
	HomeTeam     string `json:"homeTeam"`
	AwayTeam     string `json:"awayTeam"`
	WentOvertime bool   `json:"OT"`
	GWS          bool   `json:"GWS"`
}
