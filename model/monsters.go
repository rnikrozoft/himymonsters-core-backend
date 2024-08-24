package model

import (
	"github.com/google/uuid"
)

type MyMonsters struct {
	Monsters []Monster `json:"monsters,omitempty"`
}

type Monster struct {
	ID                 uuid.UUID `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	MonsterType        string    `json:"monster_type,omitempty"`
	StealChangeSuccess int       `json:"steal_change_success,omitempty"`
	KillChangeSuccess  int       `json:"kill_change_success,omitempty"`
}

type StealOrkill struct {
	FriendId  string    `json:"friend_id,omitempty"`
	MonsterId uuid.UUID `json:"monster_id,omitempty"`
}
