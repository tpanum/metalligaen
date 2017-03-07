package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

type Team struct {
	Tag   string `json:"tag"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

var teams []*Team = []*Team{
	{
		Tag:   "SE",
		Name:  "SønderjyskE",
		Color: "#fff",
	},
	{
		Tag:   "AP",
		Name:  "Aalborg Pirates",
		Color: "#C31D2B",
	},
	{
		Tag:   "EE",
		Name:  "Esbjerg Energy",
		Color: "#FEC424",
	},
	{
		Tag:   "HE",
		Name:  "Herlev Eagles",
		Color: "#000000",
	},
	{
		Tag:   "HBF",
		Name:  "Herning Blue Fox",
		Color: "#003470",
	},
	{
		Tag:   "OBD",
		Name:  "Odense Bulldogs",
		Color: "#005927",
	},
	{
		Tag:   "RMB",
		Name:  "Rødovre Mighty Bulls",
		Color: "#790002",
	},
	{
		Tag:   "OBD",
		Name:  "Odense Bulldogs",
		Color: "#005927",
	},
	{
		Tag:   "RSC",
		Name:  "Rungsted Seier Capital",
		Color: "#184FA1",
	},
	{
		Tag:   "FWH",
		Name:  "Frederikshavn White Hawks",
		Color: "#FFFFFF",
	},
	{
		Tag:   "GS",
		Name:  "Gentofte Stars",
		Color: "#FFFFFF",
	},
}

var nameToTeam map[string]*Team = make(map[string]*Team)

func init() {
	for _, t := range teams {
		nameToTeam[t.Name] = t
	}
}

func TeamTagFromName(name string) (string, error) {
	if team, ok := nameToTeam[name]; ok {
		return team.Tag, nil
	}

	return "", fmt.Errorf("Unable to find team")
}

func GetTeams(c echo.Context) error {
	return c.JSON(http.StatusOK, teams)
}
