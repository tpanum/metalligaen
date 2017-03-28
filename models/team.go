package models

import (
	"fmt"
)

type Team struct {
	Tag     string `json:"tag"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	LogoUrl string `json:"logo_url"`
}

var Teams []*Team = []*Team{
	{
		Tag:     "SE",
		Name:    "SønderjyskE",
		Color:   "#fff",
		LogoUrl: "/static/sj.png",
	},
	{
		Tag:     "AP",
		Name:    "Aalborg Pirates",
		Color:   "#C31D2B",
		LogoUrl: "/static/ap.png",
	},
	{
		Tag:     "EE",
		Name:    "Esbjerg Energy",
		Color:   "#FEC424",
		LogoUrl: "/static/ee.png",
	},
	{
		Tag:     "HE",
		Name:    "Herlev Eagles",
		Color:   "#000",
		LogoUrl: "/static/he.png",
	},
	{
		Tag:     "HBF",
		Name:    "Herning Blue Fox",
		Color:   "#003470",
		LogoUrl: "/static/hbf.png",
	},
	{
		Tag:     "OB",
		Name:    "Odense Bulldogs",
		Color:   "#005927",
		LogoUrl: "/static/ob.png",
	},
	{
		Tag:     "RMB",
		Name:    "Rødovre Mighty Bulls",
		Color:   "#790002",
		LogoUrl: "/static/rmb.png",
	},
	{
		Tag:     "RSC",
		Name:    "Rungsted Seier Capital",
		Color:   "#184FA1",
		LogoUrl: "/static/rsc.png",
	},
	{
		Tag:     "FWH",
		Name:    "Frederikshavn White Hawks",
		Color:   "#FFF",
		LogoUrl: "/static/fwh.png",
	},
	{
		Tag:     "GS",
		Name:    "Gentofte Stars",
		Color:   "#FFF",
		LogoUrl: "/static/gs.png",
	},
}

var nameToTeam map[string]*Team = make(map[string]*Team)

func init() {
	for _, t := range Teams {
		nameToTeam[t.Name] = t
	}
}

func TeamTagFromName(name string) (string, error) {
	if team, ok := nameToTeam[name]; ok {
		return team.Tag, nil
	}

	return "", fmt.Errorf("Unable to find team")
}
