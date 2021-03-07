package app

import (
	"github.com/mabaro3009/league-standings/nexus/services/standing/domain"
)

const lecURL = "https://lol.gamepedia.com/LEC/2021_Season/Spring_Season"

type currentStanding struct {
	Position int
	Name     string
}

type CurrentStandingsResp struct {
	Standings []currentStanding
}

func GetCurrentStandings(standing *domain.Standing) (*CurrentStandingsResp, error) {
	standings, err := standing.GetCurrentStandings(lecURL)
	if err != nil {
		return nil, err
	}

	standingsResp := make([]currentStanding, 0, len(standings))
	for i, st := range standings {
		currentSt := currentStanding{
			Position: i,
			Name:     st.TeamName,
		}
		standingsResp = append(standingsResp, currentSt)
	}

	currentStResp := &CurrentStandingsResp{Standings: standingsResp}

	return currentStResp, nil
}
