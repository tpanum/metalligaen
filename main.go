package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/tpanum/metalligaen/api"
	"github.com/tpanum/metalligaen/scraper"
)

func main() {
	c, _ := scraper.NewClient()
	c.GetSchedule(1254)

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/teams", api.GetTeams)
	e.GET("/matches", api.GetMatches)
	e.GET("/live", api.GetLive)

	e.Start(":8080")
}
