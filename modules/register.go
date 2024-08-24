package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/rnikrozoft/hellomymonsters-backend/model"
	"github.com/samber/lo"
)

func RpcUserRegister(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	// TODO

	var p *model.RegisterWithDeviceID
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		logger.Error(err.Error())
		return "", err
	}

	userID, _, _, err := nk.AuthenticateDevice(ctx, p.DeviceID, p.Username, false)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	w := make(map[string]int64)
	w["coins"] = 500

	if _, _, err := nk.WalletUpdate(ctx, userID, w, map[string]interface{}{}, false); err != nil {
		logger.Error(err.Error())
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	min := 1
	max := 10

	m := model.MyMonsters{
		Monsters: []model.Monster{
			{
				ID:                 lo.Must1(uuid.NewV7()),
				Name:               "monster_1",
				MonsterType:        "monster_type_1",
				StealChangeSuccess: rand.Intn(max-min+1) + min,
				KillChangeSuccess:  rand.Intn(max-min+1) + min,
			},
			{
				ID:                 lo.Must1(uuid.NewV7()),
				Name:               "monster_2",
				MonsterType:        "monster_type_1",
				StealChangeSuccess: rand.Intn(max-min+1) + min,
				KillChangeSuccess:  rand.Intn(max-min+1) + min,
			},
			{
				ID:                 lo.Must1(uuid.NewV7()),
				Name:               "monster_3",
				MonsterType:        "monster_type_1",
				StealChangeSuccess: rand.Intn(max-min+1) + min,
				KillChangeSuccess:  rand.Intn(max-min+1) + min,
			},
			{
				ID:                 lo.Must1(uuid.NewV7()),
				Name:               "monster_4",
				MonsterType:        "monster_type_1",
				StealChangeSuccess: rand.Intn(max-min+1) + min,
				KillChangeSuccess:  rand.Intn(max-min+1) + min,
			},
			{
				ID:                 lo.Must1(uuid.NewV7()),
				Name:               "monster_5",
				MonsterType:        "monster_type_1",
				StealChangeSuccess: rand.Intn(max-min+1) + min,
				KillChangeSuccess:  rand.Intn(max-min+1) + min,
			},
		},
	}

	s := []*runtime.StorageWrite{
		{
			Collection:      "Monsters",
			Key:             "MyMonsters",
			UserID:          userID,
			Value:           string(lo.Must1(json.Marshal(m))),
			PermissionWrite: 1,
			PermissionRead:  2,
		},
	}
	if _, err := nk.StorageWrite(ctx, s); err != nil {
		return "", err
	}

	return payload, nil
}
