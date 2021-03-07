package domain

import "fmt"

type standingStore interface {
}

type Standing struct {
	store standingStore
}

func NewStanding(store standingStore) *Standing {
	return &Standing{store: store}
}

func (s *Standing) Test() {
	fmt.Println("correct")
}
