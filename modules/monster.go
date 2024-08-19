package modules

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/rnikrozoft/hellomymonsters-backend/model"
)

func RpcKillMonster(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userId, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)

	if !ok || userId == "" {
		logger.Error("rpc was called by a user")
		return "", runtime.NewError("rpc is only callable via server to server", 7)
	}

	var p *model.Monster
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		logger.Error(err.Error())
		return "", err
	}

	// userID, _, _, err := nk.AuthenticateDevice(ctx, p.DeviceID, p.Username, false)
	// if err != nil {
	// 	logger.Error(err.Error())
	// 	return "", err
	// }

	// w := make(map[string]int64)
	// w["coins"] = 500

	// if _, _, err := nk.WalletUpdate(ctx, userID, w, map[string]interface{}{}, false); err != nil {
	// 	logger.Error(err.Error())
	// 	return "", err
	// }

	// m := model.MyMonsters{
	// 	Monsters: []model.Monster{
	// 		{
	// 			Name:        "Monster 1",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 2",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 3",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 4",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 5",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 6",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 7",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 8",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 9",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 10",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 11",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 		{
	// 			Name:        "Monster 12",
	// 			MonsterType: "monster_type_1",
	// 		},
	// 	},
	// }

	// s := []*runtime.StorageWrite{
	// 	{
	// 		Collection:      "Monsters",
	// 		Key:             "MyMonsters",
	// 		UserID:          userID,
	// 		Value:           string(lo.Must1(json.Marshal(m))),
	// 		PermissionWrite: 1,
	// 		PermissionRead:  2,
	// 	},
	// }
	// if _, err := nk.StorageWrite(ctx, s); err != nil {
	// 	return "", err
	// }

	return payload, nil
}
