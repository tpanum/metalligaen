package scraper_test

import (
	"testing"

	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"

	"github.com/tpanum/metalligaen/scraper"
)

func SportsadminStub() *httptest.Server {
	reqResp := map[string]string{
		"/hockeystats": fileToString(filepath.Join("..", "test", "league_page.html")),
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		for req, resp := range reqResp {
			if strings.HasPrefix(uri, req) {
				w.Write([]byte(resp))
				return
			}
		}

		http.Error(w, "not found", http.StatusNotFound)
	}))
}

func TestGetLeagueIDs(t *testing.T) {
	server := SportsadminStub()
	defer server.Close()

	c := scraper.NewTournamentClientWithDomain(server.URL)
	ids, err := c.LeagueIDsByName("Metal Ligaen")
	if err != nil {
		t.Fatalf("Unable to get league IDs: %v", err)
	}

	if len(ids) == 0 {
		t.Fatalf("Expected more than zero ids")
	}
}
