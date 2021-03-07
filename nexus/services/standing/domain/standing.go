package domain

import (
	"github.com/campoy/unique"
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
	GetRecord(team1ID models.TeamID, team2ID models.TeamID) (*models.TeamRecord, error)
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

	return s.orderStandings(unorderedStandings)
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

func (s *Standing) orderStandings(standings []*util.StandingDTO) ([]*util.StandingDTO, error) {
	wins := make(map[int][]*util.StandingDTO)
	keys := make([]int, 0)
	for _, standing := range standings {
		wins[standing.Wins] = append(wins[standing.Wins], standing)
		keys = append(keys, standing.Wins)
	}
	unique.Slice(&keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	orderedStandings := make([]*util.StandingDTO, 0, len(standings))
	for _, k := range keys {
		teams := wins[k]
		switch len(teams) {
		case 1:
			orderedStandings = append(orderedStandings, wins[k][0])
		default:
			subOrdered, err := s.solveTiebreaker(wins[k]...)
			if err != nil {
				return nil, err
			}
			orderedStandings = append(orderedStandings, subOrdered...)
		}
	}

	return orderedStandings, nil
}

// PRE: len(standings) >= 2
func (s *Standing) solveTiebreaker(standings ...*util.StandingDTO) ([]*util.StandingDTO, error) {
	orderedStandings := make([]*util.StandingDTO, 0, len(standings))
	if len(standings) == 2 {
		record, err := s.store.GetRecord(standings[0].TeamID, standings[1].TeamID)
		if err != nil {
			return nil, err
		}
		if record.Team1Wins > record.Team2Wins {
			orderedStandings = append(orderedStandings, standings[0])
			orderedStandings = append(orderedStandings, standings[1])
		} else if record.Team1Wins < record.Team2Wins {
			orderedStandings = append(orderedStandings, standings[1])
			orderedStandings = append(orderedStandings, standings[0])
		} else { // sort by second-half wins
			if standings[0].WinsSecondHalf > standings[1].WinsSecondHalf {
				orderedStandings = append(orderedStandings, standings[0])
				orderedStandings = append(orderedStandings, standings[1])
			} else if standings[0].WinsSecondHalf < standings[1].WinsSecondHalf {
				orderedStandings = append(orderedStandings, standings[1])
				orderedStandings = append(orderedStandings, standings[0])
			} else { // sort alphabetically
				if standings[0].TeamName > standings[1].TeamName {
					orderedStandings = append(orderedStandings, standings[0])
					orderedStandings = append(orderedStandings, standings[1])
				} else {
					orderedStandings = append(orderedStandings, standings[1])
					orderedStandings = append(orderedStandings, standings[0])
				}
			}
		}

		return orderedStandings, nil
	}

	tiedStandings, err := s.generateTiedStandings(standings...)
	if err != nil {
		return nil, err
	}

	return s.orderStandings(tiedStandings)
}

func (s *Standing) generateTiedStandings(standings ...*util.StandingDTO) ([]*util.StandingDTO, error) {
	tiedStandingsMap := make(map[models.TeamID]*util.StandingDTO)
	for _, standing := range standings {
		tiedStandingsMap[standing.TeamID] = &util.StandingDTO{
			TeamID:         standing.TeamID,
			TeamName:       standing.TeamName,
			Wins:           0,
			Loses:          0,
			WinsSecondHalf: standing.WinsSecondHalf,
		}
	}
	for i := 0; i < len(standings); i++ {
		for j := i + 1; j < len(standings); j++ {
			record, err := s.store.GetRecord(standings[i].TeamID, standings[j].TeamID)
			if err != nil {
				return nil, err
			}
			tiedStandingsMap[record.Team1ID].Wins += record.Team1Wins
			tiedStandingsMap[record.Team1ID].Loses += record.Team2Wins
			tiedStandingsMap[record.Team2ID].Wins += record.Team2Wins
			tiedStandingsMap[record.Team2ID].Loses += record.Team1Wins
		}
	}
	tiedStandings := make([]*util.StandingDTO, 0, len(standings))
	for _, v := range tiedStandingsMap {
		tiedStandings = append(tiedStandings, v)
	}

	return tiedStandings, nil
}

func (s *Standing) standingToDTO(standing *models.TeamStanding) (*util.StandingDTO, error) {
	name, err := s.store.GetTeamNameFromID(standing.TeamID)
	if err != nil {
		return nil, err
	}
	return &util.StandingDTO{
		TeamID:         standing.TeamID,
		TeamName:       name,
		Wins:           standing.Wins,
		Loses:          standing.Loses,
		WinsSecondHalf: standing.WinsSecondHalf,
	}, nil
}
