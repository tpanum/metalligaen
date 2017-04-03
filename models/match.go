package models

import (
	"sort"
	"strings"
	"time"

	"github.com/tpanum/metalligaen/utils"
)

type Penalty struct {
	Num      uint   `json:"num"`
	Kind     string `json:"kind"`
	Duration uint   `json:"duration"`
	Time     uint   `json:"time"`
}

type Penalties struct {
	Home []Penalty `json:"home"`
	Away []Penalty `json:"away"`
}

type Goal struct {
	ScorerNum       uint  `json:"num"`
	AssistFirstNum  *uint `json:"a1_num,omitempty"`
	AssistSecondNum *uint `json:"a2_num,omitempty"`
	Time            uint  `json:"time"`
}

type Goals struct {
	Amount  uint   `json:"amount"`
	Details []Goal `json:"details,omitempty"`
}

type Score struct {
	Home Goals `json:"home"`
	Away Goals `json:"away"`
}

type Match struct {
	ID          uint       `json:"id"`
	HomeTeamTag string     `json:"hometeam_tag"`
	AwayTeamTag string     `json:"awayteam_tag"`
	Duration    *uint      `json:"duration,omitempty"`
	Score       *Score     `json:"score,omitempty"`
	Penalties   *Penalties `json:"penalties,omitempty"`
	Spectators  *uint      `json:"spectators,omitempty"`
	TimeOfMatch time.Time  `json:"time_of_match"`
}

const DAY_FORMAT string = "Monday, 2. January 2006"

func (m *Match) UpdateDuration(dur uint) {
}

func (m Match) Day() string {
	result := m.TimeOfMatch.Format(DAY_FORMAT)

	year := time.Now().Format(" 2006")
	result = strings.TrimSuffix(result, year)

	return utils.TranslateTimeDA(result)
}

var (
	Matches   []*MatchDay
	IDToMatch map[uint]*Match = make(map[uint]*Match)
)

type MatchList []*Match

func (ml MatchList) Len() int {
	return len(ml)
}

func (ml MatchList) Less(i, j int) bool {
	return ml[i].TimeOfMatch.Before(ml[j].TimeOfMatch)
}

func (ml MatchList) Swap(i, j int) {
	ml[i], ml[j] = ml[j], ml[i]
}

type MatchDay struct {
	Day   string    `json:"day"`
	Games MatchList `json:"games"`
}

func LoadMatches(matches MatchList) {
	if len(matches) == 0 {
		return
	}

	sort.Sort(sort.Reverse(matches))

	var dayToMatches []*MatchDay
	gameDay := &MatchDay{
		Day: matches[0].Day(),
	}

	for _, m := range matches {
		IDToMatch[m.ID] = m

		day := m.Day()

		if day != gameDay.Day {
			dayToMatches = append(dayToMatches, gameDay)
			gameDay = &MatchDay{
				Day: day,
			}
		}

		gameDay.Games = append(gameDay.Games, m)
	}

	dayToMatches = append(dayToMatches, gameDay)
	Matches = dayToMatches
}
