package model

import (
	"github.com/google/uuid"
)

type MyMonsters struct {
	Monsters []Monster `json:"monsters"`
}

type Monster struct {
	ID                 uuid.UUID `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	MonsterType        string    `json:"monster_type,omitempty"`
	StealChangeSuccess int       `json:"steal_change_success,omitempty"`
	KillChangeSuccess  int       `json:"kill_change_success,omitempty"`
}

type Record struct {
	Monsters []Monster `json:"monsters"`
}

type StealOrkill struct {
	OwnerID   string    `json:"owner_id,omitempty"`
	MonsterID uuid.UUID `json:"monster_id,omitempty"`
}
