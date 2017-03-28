package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/tpanum/metalligaen/models"
	"github.com/tpanum/metalligaen/scraper"
)

func GetMatches(c echo.Context) error {
	return c.JSON(http.StatusOK, models.Matches)
}

func GetMatch(c echo.Context) error {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	m, err := scraper.MatchFromReport(uint(idInt))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, m)
}
