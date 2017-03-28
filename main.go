package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/tpanum/metalligaen/api"
	"github.com/tpanum/metalligaen/models"
	"github.com/tpanum/metalligaen/scraper"
)

func main() {

	var matches []*models.Match
	for _, id := range scraper.GetLeagueIDsByName("Metal Ligaen") {
		c, _ := scraper.NewClient(scraper.DOMAIN)
		matchesForID, _ := c.GetSchedule(id)
		matches = append(matches, matchesForID...)
	}

	models.LoadMatches(matches)

	// live, _ := scraper.NewClient(scraper.DOMAIN)
	// go live.ListenLive()

	e := echo.New()

	e.Use(middleware.Logger())

	e.Static("/static", "static/compressed")
	e.GET("/teams", api.GetTeams)
	e.GET("/matches", api.GetMatches)
	e.GET("/matches/:id", api.GetMatch)
	e.GET("/live", api.GetLive)

	e.Start(":8080")
}
