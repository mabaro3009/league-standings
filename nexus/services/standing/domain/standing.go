package domain

import (
	"github.com/mabaro3009/league-standings/nexus/models"
	"github.com/mabaro3009/league-standings/nexus/services/standing/util"
)

type standingStore interface {
	SaveTeamsInfo(names []string) error
	InitStandings() error
	UpdateRecords(result *util.ResultDTO) error
	GetTeamIDFromAbv(abv string) (models.TeamID, error)
	GetTeamStanding(teamID models.TeamID) (*models.TeamStanding, error)
	UpdateTeamStandings(team1ID, team2ID models.TeamID, result *util.ResultDTO) error
	GetAllStandings() []*models.TeamStanding
	GetTeamNameFromID(id models.TeamID) (string, error)
}

type standingScrapper interface {
	GetTeamsInfo(url string) ([]string, error)
	GetResults(url string) ([]*util.ResultScrappedDTO, error)
}

type Standing struct {
	store    standingStore
	scrapper standingScrapper
}

func NewStanding(store standingStore, scrapper standingScrapper) *Standing {
	return &Standing{store: store, scrapper: scrapper}
}

func (s *Standing) GetCurrentStandings(link string) ([]*util.StandingDTO, error) {
	if err := s.loadResults(link); err != nil {
		return nil, err
	}

	return s.getStandings()
}

func (s *Standing) loadResults(link string) error {
	names, err := s.scrapper.GetTeamsInfo(link)
	if err != nil {
		return err
	}

	if err = s.store.SaveTeamsInfo(names); err != nil {
		return err
	}

	if err = s.store.InitStandings(); err != nil {
		return err
	}

	scrappedResults, err := s.scrapper.GetResults(link)
	if err != nil {
		return err
	}

	var team1ID, team2ID models.TeamID
	var team1Standing *models.TeamStanding
	for _, scrappedResult := range scrappedResults {
		team1ID, err = s.store.GetTeamIDFromAbv(scrappedResult.Team1)
		if err != nil {
			return err
		}
		team2ID, err = s.store.GetTeamIDFromAbv(scrappedResult.Team2)
		if err != nil {
			return err
		}
		team1Standing, err = s.store.GetTeamStanding(team1ID)
		if err != nil {
			return err
		}

		result := &util.ResultDTO{
			Team1:        team1ID,
			Team2:        team2ID,
			Winner:       scrappedResult.Winner,
			IsSecondHalf: (team1Standing.Wins + team1Standing.Loses) >= util.LECGamesPerTeam/2,
		}

		if err = s.store.UpdateRecords(result); err != nil {
			return err
		}

		if err = s.store.UpdateTeamStandings(team1ID, team2ID, result); err != nil {
			return err
		}
	}

	return nil
}

func (s *Standing) getStandings() ([]*util.StandingDTO, error) {
	unorderedStandings, err := s.getUnorderedStandings()
	if err != nil {
		return nil, err
	}

	//TODO: Order standings
	return unorderedStandings, err
}

func (s *Standing) getUnorderedStandings() ([]*util.StandingDTO, error) {
	standingsFromDB := s.store.GetAllStandings()
	standings := make([]*util.StandingDTO, 0, len(standingsFromDB))
	for _, standingFromDB := range standingsFromDB {
		standing, err := s.standingToDTO(standingFromDB)
		if err != nil {
			return nil, err
		}
		standings = append(standings, standing)
	}

	return standings, nil
}

func (s *Standing) standingToDTO(standing *models.TeamStanding) (*util.StandingDTO, error) {
	name, err := s.store.GetTeamNameFromID(standing.TeamID)
	if err != nil {
		return nil, err
	}
	return &util.StandingDTO{
		TeamName:       name,
		Wins:           standing.Wins,
		Loses:          standing.Loses,
		WinsSecondHalf: standing.WinsSecondHalf,
	}, nil
}
