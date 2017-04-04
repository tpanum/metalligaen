package scraper_test

import (
	"testing"

	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"

	"github.com/tpanum/metalligaen/models"
	"github.com/tpanum/metalligaen/scraper"
)

var (
	reportsToMatches = map[string]models.Match{
		"match_report_normal.html": func() models.Match {
			var spectators uint = 517
			var duration uint = 3600
			hteam, _ := models.TeamTagFromName("Gentofte Stars")
			ateam, _ := models.TeamTagFromName("SønderjyskE")

			return models.Match{
				ID:          40463,
				Spectators:  &spectators,
				HomeTeamTag: hteam,
				AwayTeamTag: ateam,
				Duration:    &duration,
				Score: &models.Score{
					Home: models.Goals{
						Amount: 3,
						Details: []models.Goal{
							{
								ScorerNum: 50,
								Time:      374,
							},
							{
								ScorerNum: 19,
								Time:      2096,
							},
							{
								ScorerNum: 19,
								Time:      3169,
							},
						},
					},
					Away: models.Goals{
						Amount:  0,
						Details: []models.Goal{},
					},
				},
			}
		}(),
	}
)

func HockeyligaStub() *httptest.Server {
	reqResp := make(map[string]string)

	for filename, match := range reportsToMatches {
		content := fileToString(filepath.Join("..", "test", filename))

		reqResp[fmt.Sprintf("/gamesheet.aspx?gameID=%v", match.ID)] = content
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		for req, resp := range reqResp {
			if strings.HasPrefix(uri, req) {
				w.Write([]byte(resp))
				return
			}
		}

		http.Error(w, "not found", http.StatusNotFound)
	}))
}

func compareGoals(t *testing.T, expectedGoals, outputGoals []models.Goal) {
	if len(expectedGoals) != len(outputGoals) {
		t.Fatalf("Goals does not match in length (%v vs. %v)",
			len(expectedGoals),
			len(outputGoals),
		)
		return
	}

	for i, goal := range outputGoals {
		expectedGoal := expectedGoals[i]
		if expectedGoal.ScorerNum != goal.ScorerNum {
			t.Fatalf("Expected goal to be scored by \"%v\", but got \"%v\".",
				expectedGoal.ScorerNum,
				goal.ScorerNum,
			)
		}

		if expectedGoal.Time != goal.Time {
			t.Fatalf("Expected goal to be at time \"%v\", but got \"%v\".",
				expectedGoal.Time,
				goal.Time,
			)
		}

	}

}

func TestGetDetailsByID(t *testing.T) {
	server := HockeyligaStub()
	defer server.Close()

	c := scraper.NewMatchClientWithDomain(server.URL)

	for _, expectedMatch := range reportsToMatches {
		match, err := c.GetDetailsByID(expectedMatch.ID)
		if err != nil {
			t.Fatalf("Unable to get match details: " + err.Error())
		}

		if match == nil {
			t.Fatalf("Expected match not to be nil")
		}

		if match.Duration == nil && expectedMatch.Duration != nil {
			t.Fatalf("Expected match duration not to be nil")
		}

		if *match.Duration != *expectedMatch.Duration {
			t.Fatalf("Expected match duration to be %v, but got %v.",
				expectedMatch.Duration,
				match.Duration,
			)
		}

		if match.Spectators == nil && expectedMatch.Spectators != nil {
			t.Fatalf("Expected match spectators not to be nil")
		}

		if *match.Spectators != *expectedMatch.Spectators {
			t.Fatalf("Expected match spectators to be %v, but got %v.",
				expectedMatch.Spectators,
				match.Spectators,
			)
		}

		if match.HomeTeamTag != expectedMatch.HomeTeamTag {
			t.Fatalf("Expected match home team tag to be \"%v\", but got \"%v\".",
				expectedMatch.HomeTeamTag,
				match.HomeTeamTag,
			)
		}

		if match.AwayTeamTag != expectedMatch.AwayTeamTag {
			t.Fatalf("Expected match away team tag to be \"%v\", but got \"%v\".",
				expectedMatch.AwayTeamTag,
				match.AwayTeamTag,
			)
		}

		if match.Score == nil && expectedMatch.Score != nil {
			t.Fatalf("Expected match score not to be nil")
		}

		if match.Score.Home.Amount != expectedMatch.Score.Home.Amount {
			t.Fatalf("Expected match home score to be \"%v\", but got \"%v\".",
				expectedMatch.Score.Home.Amount,
				match.Score.Home.Amount,
			)
		}

		if match.Score.Away.Amount != expectedMatch.Score.Away.Amount {
			t.Fatalf("Expected match away score to be \"%v\", but got \"%v\".",
				expectedMatch.Score.Away.Amount,
				match.Score.Away.Amount,
			)
		}

		if expectedMatch.Score.Home.Details != nil {
			compareGoals(t,
				expectedMatch.Score.Home.Details,
				match.Score.Home.Details,
			)
		}

		if expectedMatch.Score.Home.Details != nil {
			compareGoals(t,
				expectedMatch.Score.Away.Details,
				match.Score.Away.Details,
			)
		}
	}
}
