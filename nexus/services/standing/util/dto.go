package util

import "github.com/mabaro3009/league-standings/nexus/models"

type ResultScrappedDTO struct {
	Team1  string
	Team2  string
	Winner bool // false = team1, true = team2
}

type ResultDTO struct {
	Team1        models.TeamID
	Team2        models.TeamID
	Winner       bool // false = team1, true = team2
	IsSecondHalf bool
}

type StandingDTO struct {
	TeamID         models.TeamID
	TeamName       string
	Wins           int
	Loses          int
	WinsSecondHalf int
}
