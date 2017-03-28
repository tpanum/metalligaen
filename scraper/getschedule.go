package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"encoding/json"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"

	"github.com/tpanum/metalligaen/models"
)

func scoreFromGameResult(result string) (*models.Score, error) {
	scoreSplit := strings.Split(result, " - ")
	var score *models.Score

	if len(scoreSplit) == 2 {
		homeScoreStr, awayScoreStr := scoreSplit[0], scoreSplit[1]

		val, err := strconv.Atoi(homeScoreStr)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse home score to integer")
		}

		score = &models.Score{
			Home: models.Goals{
				Amount: uint(val),
			},
		}

		val, err = strconv.Atoi(awayScoreStr)
		if err != nil {
			return nil, fmt.Errorf("Unable to parse away score to integer")
		}

		score.Away = models.Goals{
			Amount: uint(val),
		}
	}

	return score, nil
}

func (c *Client) GetSchedule(id int) ([]*models.Match, error) {
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
		return nil, err
	}

	output := c.HookEvent("schedule")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	rawData := <-output

	var jsonData struct {
		Games []rawMatch `json:"schedule"`
	}

	if err := json.Unmarshal(rawData, &jsonData); err != nil {
		return nil, err
	}

	matches := make([]*models.Match, len(jsonData.Games))
	for i, m := range jsonData.Games {
		timeOfMatch, err := time.Parse("2006-01-02T15:04", m.Date[0:11]+m.Time)
		if err != nil {
			fmt.Println(err)
		}

		hometeam, err := models.TeamTagFromName(m.HomeTeam)
		if err != nil {
			continue
		}

		awayteam, err := models.TeamTagFromName(m.AwayTeam)
		if err != nil {
			continue
		}

		score, err := scoreFromGameResult(m.GameScore)

		matches[i] = &models.Match{
			ID:          uint(m.ID),
			HomeTeamTag: hometeam,
			AwayTeamTag: awayteam,
			Score:       score,
			TimeOfMatch: timeOfMatch,
		}
	}

	return matches, nil
}

var (
	TOURNEYID_REGEXP = regexp.MustCompile(`tournamentid=([0-9]+)`)
)

func GetLeagueIDsByName(name string) []int {
	resp, err := http.Get("http://www.sportsadmin.dk/hockeystats/")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil
	}

	var result []int
	doc.Find("tr").Each(func(i int, row *goquery.Selection) {
		cols := row.Find("td")

		tournamentName := cols.First().Text()
		if strings.Contains(tournamentName, name) {
			link, ok := row.Find("a").Last().Attr("href")
			if !ok {
				return
			}

			matches := TOURNEYID_REGEXP.FindStringSubmatch(link)
			if len(matches) > 0 {
				idStr := matches[1]
				id, _ := strconv.Atoi(idStr)
				result = append(result, id)
			}
		}
	})

	return result
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
