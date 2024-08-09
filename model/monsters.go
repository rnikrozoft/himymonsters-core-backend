package model

type MyMonsters struct {
	Monsters []Monster `json:"monsters,omitempty"`
}

type Monster struct {
	Name        string `json:"name,omitempty"`
	Health      int    `json:"health,omitempty"`
	Attack      int    `json:"attack,omitempty"`
	Defence     int    `json:"defence,omitempty"`
	Experience  int    `json:"experience,omitempty"`
	Level       int    `json:"level,omitempty"`
	MonsterType string `json:"monster_type,omitempty"`
}
