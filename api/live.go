package api

import (
	"github.com/labstack/echo"

	"github.com/tpanum/metalligaen/models"
	"golang.org/x/net/websocket"
)

var (
	liveMatches = map[int]models.Match{}
)

var LivePool = NewWSPool()

func init() {
	go LivePool.Run()
}

func GetLive(c echo.Context) error {
	websocket.Handler(func(conn *websocket.Conn) {
		LivePool.register <- conn
		defer conn.Close()

		websocket.Message.Send(conn, "Hey guys")

		for {

		}
	}).ServeHTTP(c.Response(), c.Request())

	return nil
}
