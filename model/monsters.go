package model

import (
	"github.com/google/uuid"
)

type MyMonsters struct {
	Monsters []Monster `json:"monsters"`
}

type Monster struct {
	ID          uuid.UUID           `json:"id,omitempty"`
	Name        string              `json:"name,omitempty"`
	MonsterType string              `json:"monster_type,omitempty"`
	Steal       StealOrKillSettings `json:"steal,omitempty"`
	Kill        StealOrKillSettings `json:"kill,omitempty"`
}

type Record struct {
	Monsters []Monster `json:"monsters"`
}

type StealOrkill struct {
	OwnerID   string    `json:"owner_id,omitempty"`
	MonsterID uuid.UUID `json:"monster_id,omitempty"`
}

type StealOrKillSettings struct {
	RateSuccess int   `json:"rate_success"`
	Price       int64 `json:"price"`
}
