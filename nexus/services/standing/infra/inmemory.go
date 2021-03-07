package infra

import "github.com/mabaro3009/league-standings/nexus/models"

type StandingStoreInMemory struct {
	Leagues   map[models.LeagueID]*models.League
	Teams     map[models.TeamID]*models.Team
	Standings map[models.TeamStandingID]*models.TeamStanding
	Records   map[models.TeamRecordID]*models.TeamRecord
}

func NewStandingStoreInMemory() *StandingStoreInMemory {
	return &StandingStoreInMemory{
		Leagues:   make(map[models.LeagueID]*models.League),
		Teams:     make(map[models.TeamID]*models.Team),
		Standings: make(map[models.TeamStandingID]*models.TeamStanding),
		Records:   make(map[models.TeamRecordID]*models.TeamRecord),
	}
}
