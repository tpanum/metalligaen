package api

import (
	"time"

	"github.com/labstack/echo"
	"net/http"
)

type penalty struct {
	Num      uint
	Duration uint
	Time     uint
}

type Situation struct {
	Goals     uint      `json:"goals"`
	Penalties []penalty `json:"penalties,omitempty"`
}

type Score struct {
	Home Situation `json:"home"`
	Away Situation `json:"away"`
}

type Match struct {
	ID          uint      `json:"id"`
	HomeTeamTag string    `json:"hometeam_tag"`
	AwayTeamTag string    `json:"awayteam_tag"`
	Score       *Score    `json:"score,omitempty"`
	TimeOfMatch time.Time `json:"time_of_match"`
}

var matches []*Match

func LoadMatches(m []*Match) {
	matches = m
}

func GetMatches(c echo.Context) error {
	return c.JSON(http.StatusOK, matches)
}
