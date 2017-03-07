package api

import (
	"time"

	"github.com/labstack/echo"
	"net/http"
)

type Match struct {
	ID          uint      `json:"id"`
	HomeTeamTag string    `json:"hometeam_tag"`
	AwayTeamTag string    `json:"awayteam_tag"`
	TimeOfMatch time.Time `json:"time_of_match"`
}

var matches []*Match

func LoadMatches(m []*Match) {
	matches = m
}

func GetMatches(c echo.Context) error {
	return c.JSON(http.StatusOK, matches)
}
