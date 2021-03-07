package models

import (
	"github.com/nu7hatch/gouuid"
)

type TeamID string

func NewTeamID() TeamID {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return TeamID(id.String())
}

type Team struct {
	ID   TeamID `json:"id"`
	Name string `json:"name"`
	Abv  string `json:"abv"`
}

type TeamStandingID string

func NewTeamStandingID() TeamStandingID {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return TeamStandingID(id.String())
}

type TeamStanding struct {
	ID             TeamStandingID `json:"id"`
	TeamID         TeamID         `json:"team_id"`
	Wins           int            `json:"wins"`
	Loses          int            `json:"loses"`
	WinsSecondHalf int            `json:"wins_second_half"`
}

type TeamRecordID string

func NewTeamRecordID() TeamRecordID {
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	return TeamRecordID(id.String())
}

type TeamRecord struct {
	ID        TeamRecordID `json:"id"`
	Team1ID   TeamID       `json:"team1_id"`
	Team2ID   TeamID       `json:"team2_id"`
	Team1Wins int          `json:"team1_wins"`
	Team2Wins int          `json:"team2_wins"`
}
