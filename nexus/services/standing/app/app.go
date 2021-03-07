package app

import (
	"github.com/mabaro3009/league-standings/nexus/services/standing/domain"
	"github.com/mabaro3009/league-standings/nexus/services/standing/util"
)

const lecURL = "https://lol.gamepedia.com/LEC/2021_Season/Spring_Season"

func GetCurrentStandings(standing *domain.Standing) ([]*util.StandingDTO, error) {
	return standing.GetCurrentStandings(lecURL)
}
