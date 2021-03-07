package standing

import (
	"github.com/mabaro3009/league-standings/nexus/services/standing/app"
	"github.com/mabaro3009/league-standings/nexus/services/standing/domain"
	"github.com/mabaro3009/league-standings/nexus/services/standing/infra"
)

type Service struct {
	standing *domain.Standing
}

func New() *Service {
	return &Service{standing: domain.NewStanding(infra.NewStandingStoreInMemory(), &infra.Scrapper{})}
}

func (s *Service) GetCurrentStandings() *app.CurrentStandingsResp {
	standings, err := app.GetCurrentStandings(s.standing)
	if err != nil {
		panic(err)
	}
	return standings
}
