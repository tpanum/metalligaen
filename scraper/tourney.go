package scraper

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var (
	TOURNEYID_REGEXP = regexp.MustCompile(`tournamentid=([0-9]+)`)
)

type TournamentClient interface {
	LeagueIDsByName(string) ([]int, error)
}

func NewTournamentClient() TournamentClient {
	return &sportsAdminClient{
		domain: "http://hockeyligaen.dk",
	}
}

func NewTournamentClientWithDomain(domain string) TournamentClient {
	return &sportsAdminClient{
		domain: domain,
	}
}

type sportsAdminClient struct {
	domain string
}

func (sac *sportsAdminClient) LeagueIDsByName(name string) ([]int, error) {
	resp, err := http.Get(fmt.Sprintf("%s/hockeystats/", sac.domain))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return nil, err
	}

	var result []int
	doc.Find("tr").Each(func(i int, row *goquery.Selection) {
		cols := row.Find("td")

		tournamentName := cols.First().Text()
		if strings.Contains(tournamentName, name) {
			link, ok := row.Find("a").Last().Attr("href")
			if !ok {
				return
			}

			matches := TOURNEYID_REGEXP.FindStringSubmatch(link)
			if len(matches) > 0 {
				idStr := matches[1]
				id, _ := strconv.Atoi(idStr)
				result = append(result, id)
			}
		}
	})

	return result, nil
}
