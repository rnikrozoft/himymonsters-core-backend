package modules

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/rnikrozoft/hellomymonsters-backend/model"
)

func RpcStealMonster(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	var p *model.StealOrkill
	json.Unmarshal([]byte(payload), &p)

	objectIds := []*runtime.StorageRead{
		{
			Collection: "Monsters",
			Key:        "MyMonsters",
			UserID:     p.FriendId,
		},
	}

	records, err := nk.StorageRead(ctx, objectIds)
	if err != nil {
		logger.WithField("err", err).Error("Storage read error.")
	} else {
		m := &model.Monster{}
		json.Unmarshal([]byte(records[0].Value), m)
	}

	// objectIds := []*runtime.StorageDelete{
	// 	{
	// 		Collection: "Monsters",
	// 		Key:        "MyMonsters",
	// 		UserID:     p.FriendId,
	// 	},
	// }

	// err := nk.StorageDelete(ctx, objectIds)
	// if err != nil {
	// 	logger.WithField("err", err).Error("Storage delete error.")
	// }
	return payload, nil
}
