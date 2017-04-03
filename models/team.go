package models

import (
	"fmt"
)

type colorTheme struct {
	Main      string `json:"main"`
	Secondary string `json:"secondary"`
}

type Team struct {
	Tag     string     `json:"tag"`
	Name    string     `json:"name"`
	Colors  colorTheme `json:"colors"`
	LogoUrl string     `json:"logo_url"`
}

var Teams []*Team = []*Team{
	{
		Tag:  "SE",
		Name: "SønderjyskE",
		Colors: colorTheme{
			Main:      "#fff",
			Secondary: "#344480",
		},
		LogoUrl: "/static/sj.png",
	},
	{
		Tag:  "AP",
		Name: "Aalborg Pirates",
		Colors: colorTheme{
			Main:      "#C31D2B",
			Secondary: "#182336",
		},
		LogoUrl: "/static/ap.png",
	},
	{
		Tag:  "EE",
		Name: "Esbjerg Energy",
		Colors: colorTheme{
			Main:      "#FEC424",
			Secondary: "#10213E",
		},
		LogoUrl: "/static/ee.png",
	},
	{
		Tag:  "HE",
		Name: "Herlev Eagles",
		Colors: colorTheme{
			Main:      "#000",
			Secondary: "#FEC931",
		},
		LogoUrl: "/static/he.png",
	},
	{
		Tag:  "HBF",
		Name: "Herning Blue Fox",
		Colors: colorTheme{
			Main:      "#003470",
			Secondary: "#FABA00",
		},
		LogoUrl: "/static/hbf.png",
	},
	{
		Tag:  "OB",
		Name: "Odense Bulldogs",
		Colors: colorTheme{
			Main:      "#005927",
			Secondary: "#000",
		},
		LogoUrl: "/static/ob.png",
	},
	{
		Tag:  "RMB",
		Name: "Rødovre Mighty Bulls",
		Colors: colorTheme{
			Main:      "#790002",
			Secondary: "#000",
		},
		LogoUrl: "/static/rmb.png",
	},
	{
		Tag:  "RSC",
		Name: "Rungsted Seier Capital",
		Colors: colorTheme{
			Main:      "#184FA1",
			Secondary: "#D3253C",
		},
		LogoUrl: "/static/rsc.png",
	},
	{
		Tag:  "FWH",
		Name: "Frederikshavn White Hawks",
		Colors: colorTheme{
			Main:      "#FFF",
			Secondary: "#F8A536",
		},
		LogoUrl: "/static/fwh.png",
	},
	{
		Tag:  "GS",
		Name: "Gentofte Stars",
		Colors: colorTheme{
			Main:      "#FFF",
			Secondary: "#002B60",
		},
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
