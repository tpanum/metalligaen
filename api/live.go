package api

import (
	"github.com/labstack/echo"

	"golang.org/x/net/websocket"
)

type Penalty struct {
	FullName     string
	Number       uint
	DurationLeft uint
}

type MatchLiveDetails struct {
	MatchID         string `json:"match_id"`
	HomeTeamScore   uint   `json:"hometeam_score"`
	AwayTeamScore   uint   `json:"awayteam_score"`
	GameTimeSeconds uint   `json:"gametime_seconds"`
	Penalties       struct {
		Home []Penalty `json:"home"`
		Away []Penalty `json:"away"`
	} `json:"penalties"`
}

var LivePool = NewWSPool()

func init() {
	go LivePool.Run()
}

func GetLive(c echo.Context) error {
	websocket.Handler(func(conn *websocket.Conn) {
		LivePool.register <- conn
		defer conn.Close()

		for {

		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
