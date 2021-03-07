package models

import (
	"github.com/nu7hatch/gouuid"
)

type LeagueID uuid.UUID
type League struct {
	ID                LeagueID `json:"id"`
	Name              string   `json:"name"`
	GamesBetweenTeams int      `json:"games_between_teams"`
	BestOf            int      `json:"best_of"`
}

type TeamID uuid.UUID
type Team struct {
	ID       TeamID   `json:"id"`
	Name     string   `json:"name"`
	Abv      string   `json:"abv"`
	LeagueID LeagueID `json:"league_id"`
}

type TeamStandingID uuid.UUID
type TeamStanding struct {
	ID             TeamStandingID `json:"id"`
	TeamID         TeamID         `json:"team_id"`
	Wins           int            `json:"wins"`
	Loses          int            `json:"loses"`
	WinsSecondHalf int            `json:"wins_second_half"`
}

type TeamRecordID uuid.UUID
type TeamRecord struct {
	ID        TeamRecordID `json:"id"`
	Team1ID   TeamID       `json:"team1_id"`
	Team2ID   TeamID       `json:"team2_id"`
	Team1Wins int          `json:"team1_wins"`
	Team2Wins int          `json:"team2_wins"`
}
