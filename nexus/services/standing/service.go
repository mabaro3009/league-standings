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
	return &Service{standing: domain.NewStanding(infra.NewStandingStoreInMemory())}
}

func (s *Service) Test() {
	app.Test(s.standing)
}
