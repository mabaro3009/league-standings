package main

import (
	"fmt"
	"github.com/mabaro3009/league-standings/nexus/services/standing"
)

func main() {
	s := standing.New()
	standings := s.GetCurrentStandings()
	for _, stand := range standings {
		fmt.Println(stand)
	}
}
