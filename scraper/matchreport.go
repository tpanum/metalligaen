package scraper

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/tpanum/metalligaen/models"
	"github.com/tpanum/metalligaen/utils"
)

type MatchClient interface {
	GetDetailsByID(uint) (*models.Match, error)
}

func NewMatchClient() MatchClient {
	return &hockeyLigaClient{
		domain: "http://hockeyligaen.dk",
	}
}

func NewMatchClientWithDomain(domain string) MatchClient {
	return &hockeyLigaClient{
		domain: domain,
	}
}

type hockeyLigaClient struct {
	domain string
}

func (hlc *hockeyLigaClient) GetDetailsByID(id uint) (*models.Match, error) {
	resp, err := http.Get(fmt.Sprintf("%s/gamesheet.aspx?gameID=%v", hlc.domain, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	homePenalties, awayPenalties, err := penaltiesFromReport(doc)
	if err != nil {
		return nil, err
	}

	spectators, err := spectatorsFromReport(doc)
	if err != nil {
		return nil, err
	}

	scoreInfo, err := durationFromReport(doc)
	if err != nil {
		return nil, err
	}

	score, err := scoreFromReport(doc)
	if err != nil {
		return nil, err
	}

	homeTeam, awayTeam, err := teamsFromReport(doc)
	if err != nil {
		return nil, err
	}

	homeTeamTag, err := models.TeamTagFromName(homeTeam)
	if err != nil {
		return nil, err
	}

	awayTeamTag, err := models.TeamTagFromName(awayTeam)
	if err != nil {
		return nil, err
	}

	return &models.Match{
		ID:          id,
		HomeTeamTag: homeTeamTag,
		AwayTeamTag: awayTeamTag,
		Duration:    &scoreInfo.Duration,
		Score:       score,
		Penalties: &models.Penalties{
			Home: homePenalties,
			Away: awayPenalties,
		},
		Spectators: &spectators,
	}, nil
}

func scoreFromReport(doc *goquery.Document) (*models.Score, error) {
	ids := []string{
		"ctl00_ContentPlaceHolder1_lblGoalsHome",
		"ctl00_ContentPlaceHolder1_lblGoalsAway",
	}

	score := &models.Score{}

	for j, id := range ids {
		teamGoals := doc.Find("#" + id).Parent().Find("tr")

		for i := 2; i < teamGoals.Length(); i++ {
			goal := teamGoals.Eq(i)
			fields := goal.Find("td")

			goalNum := strings.TrimSpace(fields.Eq(0).Text())
			if goalNum == "" {
				break
			}

			goalTime := strings.TrimSpace(fields.Eq(1).Text())
			time, err := utils.GameTimeToSeconds(goalTime)
			if err != nil {
				return nil, err
			}

			g := models.Goal{
				Time: time,
			}

			scorer, err := strconv.Atoi(strings.TrimSpace(fields.Eq(2).Text()))
			if err != nil {
				return nil, err
			}
			g.ScorerNum = uint(scorer)

			firstAssister, err := strconv.Atoi(strings.TrimSpace(fields.Eq(3).Text()))
			if err == nil {
				firstAssisterUint := uint(firstAssister)
				g.AssistFirstNum = &firstAssisterUint
			}

			secondAssister, err := strconv.Atoi(strings.TrimSpace(fields.Eq(4).Text()))
			if err == nil {
				secondAssisterUint := uint(secondAssister)
				g.AssistSecondNum = &secondAssisterUint
			}

			switch j {
			case 0:
				score.Home.Details = append(score.Home.Details, g)
			case 1:
				score.Away.Details = append(score.Away.Details, g)
			}
		}
	}

	score.Home.Amount = uint(len(score.Home.Details))
	score.Away.Amount = uint(len(score.Away.Details))

	return score, nil
}

func spectatorsFromReport(doc *goquery.Document) (uint, error) {
	spectators, err := strconv.Atoi(
		strings.TrimSpace(
			doc.Find("#ctl00_ContentPlaceHolder1_lblSpectators").Text(),
		),
	)

	if err != nil {
		return 0, err
	}

	return uint(spectators), nil
}

func teamsFromReport(doc *goquery.Document) (string, string, error) {
	homeTeam := doc.Find("#ctl00_ContentPlaceHolder1_lblHjemmeHold").Text()
	awayTeam := doc.Find("#ctl00_ContentPlaceHolder1_lblUdeHold").Text()

	return homeTeam, awayTeam, nil
}

func penaltiesFromReport(doc *goquery.Document) ([]models.Penalty, []models.Penalty, error) {
	ids := []string{
		"ctl00_ContentPlaceHolder1_lblPenaltiesHome",
		"ctl00_ContentPlaceHolder1_lblPenaltiesAway",
	}

	homePenalties, awayPenalties := []models.Penalty{}, []models.Penalty{}

	for j, id := range ids {
		penalties := doc.Find("#" + id).Parent()

		teamPenalties := penalties.Find("tr")
		for i := 2; i < teamPenalties.Length(); i++ {
			penalty := teamPenalties.Eq(i)
			fields := penalty.Find("td")

			time := strings.TrimSpace(fields.Eq(0).Text())
			if time == "" {
				break
			}

			eventTime, err := utils.GameTimeToSeconds(time)
			if err != nil {
				continue
			}

			penaltier, _ := strconv.Atoi(strings.TrimSpace(fields.Eq(1).Text()))
			penaltyTimeMins, _ := strconv.Atoi(strings.TrimSpace(fields.Eq(2).Text()))
			offence := strings.TrimSpace(fields.Eq(3).Text())

			p := models.Penalty{
				Num:      uint(penaltier),
				Kind:     offence,
				Duration: uint(penaltyTimeMins * 60),
				Time:     eventTime,
			}

			switch j {
			case 0:
				homePenalties = append(homePenalties, p)
			case 1:
				awayPenalties = append(awayPenalties, p)
			}
		}
	}

	return homePenalties, awayPenalties, nil
}

type scoreInfo struct {
	Duration uint
}

var (
	PARSE_GOALS_ERR = fmt.Errorf("Unable to parse goals")
)

func durationFromReport(doc *goquery.Document) (*scoreInfo, error) {
	tables, err := tablesFromHeader(doc, "Goals")
	if err != nil {
		return nil, err
	}

	if len(tables) < 2 {
		return nil, PARSE_GOALS_ERR
	}

	gameState := ""
	var maxGoalDur uint

	for _, team := range tables {
		for _, goal := range team {
			dur, err := utils.GameTimeToSeconds(goal["Time"])
			if err != nil {
				return nil, PARSE_GOALS_ERR
			}

			if dur > maxGoalDur {
				maxGoalDur = dur
				gameState = goal["GS"]
			}
		}
	}

	var gameDuration uint = 3600

	if maxGoalDur > 3600 && gameState == "GWS" {
		gameDuration = 4500
	}

	if maxGoalDur > 3600 && gameState != "GWS" {
		gameDuration = maxGoalDur
	}

	return &scoreInfo{
		Duration: gameDuration,
	}, nil
}

func tablesFromHeader(report *goquery.Document, header string) ([][]map[string]string, error) {
	var result [][]map[string]string

	var err error
	report.Find(`b:contains("` + header + `")`).EachWithBreak(func(i int, s *goquery.Selection) bool {
		if strings.TrimSpace(s.Text()) != header {
			return true
		}

		table := s.Parent().Parent().Parent()
		rows := table.Find("tr")

		if rows.Length() < 3 {
			err = fmt.Errorf("Not enough rows (table \"%s\"seems to be empty?)", header)
			return false
		}

		var keys []string
		headerFields := rows.Eq(1).Find("td")
		for i := 0; i < headerFields.Length(); i++ {
			hfield := strings.TrimSpace(headerFields.Eq(i).Text())

			keys = append(keys, hfield)
		}

		var mapRows []map[string]string
		for i := 2; i < rows.Length(); i++ {
			row := rows.Eq(i)
			fields := row.Find("td")

			mapRow := make(map[string]string)
			nFields := fields.Length()
			blanks := 0

			for j := 0; j < nFields; j++ {
				val := strings.TrimSpace(fields.Eq(j).Text())
				if val == "" {
					blanks += 1
				}

				mapRow[keys[j]] = val
			}

			if blanks == nFields {
				break
			}

			mapRows = append(mapRows, mapRow)
		}

		result = append(result, mapRows)

		return true
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
