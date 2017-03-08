package scraper

import (
	"fmt"
	"strconv"
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
			continue
		}

		awayteam, err := api.TeamTagFromName(m.AwayTeam)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var score *api.Score
		scoreSplit := strings.Split(m.GameScore, " - ")
		if len(scoreSplit) != 0 {
			homeScoreStr, awayScoreStr := scoreSplit[0], scoreSplit[1]

			val, err := strconv.Atoi(homeScoreStr)
			if err != nil {
				fmt.Println("Unable to convert home score to int")
				continue
			}

			score = &api.Score{
				Home: api.Situation{
					Goals: uint(val),
				},
			}

			val, err = strconv.Atoi(awayScoreStr)
			if err != nil {
				fmt.Println("Unable to convert away score to int")
				continue
			}

			score.Away = api.Situation{
				Goals: uint(val),
			}

		}

		matches[i] = &api.Match{
			ID:          uint(m.ID),
			HomeTeamTag: hometeam,
			AwayTeamTag: awayteam,
			Score:       score,
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
	GameScore    string `json:"gameResult"`
	WentOvertime bool   `json:"OT"`
	GWS          bool   `json:"GWS"`
}
