package infra

import (
	"errors"
	"github.com/antchfx/htmlquery"
	"github.com/mabaro3009/league-standings/nexus/services/standing/util"
)

type Scrapper struct{}

func (s *Scrapper) GetTeamsInfo(url string) ([]string, error) {
	names := make([]string, 0)
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}

	results := htmlquery.Find(doc, "//th[@class='tournament-roster-header']/div/a")
	if len(results) == 0 {
		return nil, errors.New("no teams found")
	}

	for _, result := range results {
		names = append(names, htmlquery.InnerText(result))
	}

	return names, nil
}

func (s *Scrapper) GetResults(url string) ([]*util.ResultScrappedDTO, error) {
	doc, err := htmlquery.LoadURL(url)
	if err != nil {
		return nil, err
	}

	teams := htmlquery.Find(doc, "//table[@class='wikitable md-table']/tbody/tr/td/span[@class='team']/span[@class='teamname']")
	if len(teams) == 0 {
		return nil, errors.New("no results found")
	}

	results := htmlquery.Find(doc, "//table[@class='wikitable md-table']/tbody/tr/td")
	if len(results) == 0 {
		return nil, errors.New("no results found")
	}

	resultsDTO := make([]*util.ResultScrappedDTO, 0)
	var i int
	for _, result := range results {
		resultStr := htmlquery.InnerText(result)
		if resultStr == "1 - 0" || resultStr == "0 - 1" {
			resultDTO := &util.ResultScrappedDTO{
				Team1:  htmlquery.InnerText(teams[i*2]),
				Team2:  htmlquery.InnerText(teams[i*2+1]),
				Winner: resultStr == "0 - 1",
			}
			resultsDTO = append(resultsDTO, resultDTO)
			i++
		}
	}

	return resultsDTO, nil
}
