package infra

import (
	"errors"
	"github.com/mabaro3009/league-standings/nexus/models"
	"github.com/mabaro3009/league-standings/nexus/services/standing/util"
)

var ErrNotFound = errors.New("not found")

type StandingStoreInMemory struct {
	teams     map[models.TeamID]*models.Team
	standings map[models.TeamStandingID]*models.TeamStanding
	records   map[models.TeamRecordID]*models.TeamRecord
}

func NewStandingStoreInMemory() *StandingStoreInMemory {
	return &StandingStoreInMemory{
		teams:     make(map[models.TeamID]*models.Team),
		standings: make(map[models.TeamStandingID]*models.TeamStanding),
		records:   make(map[models.TeamRecordID]*models.TeamRecord),
	}
}

func (s *StandingStoreInMemory) SaveTeamsInfo(names []string) error {
	for _, name := range names {
		team := &models.Team{
			ID:   models.NewTeamID(),
			Name: name,
			Abv:  util.LECTeamsNametoAbv[name],
		}
		s.teams[team.ID] = team
	}

	return nil
}

func (s *StandingStoreInMemory) InitStandings() error {
	for _, team := range s.teams {
		standing := &models.TeamStanding{
			ID:             models.NewTeamStandingID(),
			TeamID:         team.ID,
			Wins:           0,
			Loses:          0,
			WinsSecondHalf: 0,
		}
		s.standings[standing.ID] = standing
	}

	return nil
}

func (s *StandingStoreInMemory) loadRecords(result *util.ResultDTO) (*models.TeamRecord, *models.TeamRecord, error) {
	var record1, record2 *models.TeamRecord
	for _, v := range s.records {
		if v.Team1ID == result.Team1 && v.Team2ID == result.Team2 {
			record1 = v
		}
		if v.Team1ID == result.Team2 && v.Team2ID == result.Team1 {
			record2 = v
		}
	}
	if record1 == nil || record2 == nil {
		return record1, record2, ErrNotFound
	}

	return record1, record2, nil
}

func (s *StandingStoreInMemory) UpdateRecords(result *util.ResultDTO) error {
	record1, record2, err := s.loadRecords(result)
	if err == ErrNotFound {
		record1 = &models.TeamRecord{
			ID:      models.NewTeamRecordID(),
			Team1ID: result.Team1,
			Team2ID: result.Team2,
		}
		record2 = &models.TeamRecord{
			ID:      models.NewTeamRecordID(),
			Team1ID: result.Team2,
			Team2ID: result.Team1,
		}
		if result.Winner {
			record1.Team2Wins = 1
			record2.Team1Wins = 1
		}
		s.records[record1.ID] = record1
		s.records[record2.ID] = record2

		return nil
	}
	if err != nil {
		return err
	}

	if result.Winner {
		record1.Team2Wins++
		record2.Team1Wins++
	} else {
		record1.Team1Wins++
		record2.Team2Wins++
	}

	return nil
}

func (s *StandingStoreInMemory) GetTeamIDFromAbv(abv string) (models.TeamID, error) {
	for _, v := range s.teams {
		if v.Abv == abv {
			return v.ID, nil
		}
	}

	return "", ErrNotFound
}

func (s *StandingStoreInMemory) GetTeamStanding(teamID models.TeamID) (*models.TeamStanding, error) {
	for _, v := range s.standings {
		if v.TeamID == teamID {
			return v, nil
		}
	}

	return nil, ErrNotFound
}

func (s *StandingStoreInMemory) UpdateTeamStandings(team1ID, team2ID models.TeamID, result *util.ResultDTO) error {
	team1Standing, err := s.GetTeamStanding(team1ID)
	if err != nil {
		return err
	}
	team2Standing, err := s.GetTeamStanding(team2ID)
	if err != nil {
		return err
	}

	if result.Winner {
		team1Standing.Loses++
		team2Standing.Wins++
		if result.IsSecondHalf {
			team2Standing.WinsSecondHalf++
		}
	} else {
		team2Standing.Loses++
		team1Standing.Wins++
		if result.IsSecondHalf {
			team1Standing.WinsSecondHalf++
		}
	}

	return nil
}

func (s *StandingStoreInMemory) GetAllStandings() []*models.TeamStanding {
	standings := make([]*models.TeamStanding, 0, len(s.standings))
	for _, v := range s.standings {
		standings = append(standings, v)
	}

	return standings
}

func (s *StandingStoreInMemory) GetTeamNameFromID(id models.TeamID) (string, error) {
	team, ok := s.teams[id]
	if !ok {
		return "", ErrNotFound
	}

	return team.Name, nil
}
