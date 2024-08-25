package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/mroth/weightedrand/v2"
	"github.com/rnikrozoft/hellomymonsters-backend/model"
	"github.com/samber/lo"
)

func RpcKillMonster(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	var p *model.StealOrkill
	json.Unmarshal([]byte(payload), &p)

	objectIds := []*runtime.StorageRead{
		{
			Collection: "Monsters",
			Key:        "MyMonsters",
			UserID:     p.OwnerID,
		},
	}

	killed := "false"
	records, err := nk.StorageRead(ctx, objectIds)
	if err != nil {
		logger.WithField("err", err).Error("Storage read error.")
	} else {

		var record model.Record
		err := json.Unmarshal([]byte(records[0].Value), &record)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
		}

		monsters := lo.Associate(record.Monsters, func(f model.Monster) (uuid.UUID, model.Monster) {
			return f.ID, f
		})

		monster := monsters[p.MonsterID]
		max := 10 - monster.KillChangeSuccess
		if max <= 0 {
			delete(monsters, p.MonsterID)
			killed = "true"
		} else {
			chooser, _ := weightedrand.NewChooser(
				weightedrand.NewChoice(true, monster.KillChangeSuccess),
				weightedrand.NewChoice(false, max),
			)
			result := chooser.Pick()

			if result {
				delete(monsters, p.MonsterID)
				killed = "true"
			}
		}

		if killed == "true" {
			value := []model.Monster{}
			for _, v := range monsters {
				new := model.Monster{
					ID:                 v.ID,
					Name:               v.Name,
					MonsterType:        v.MonsterType,
					StealChangeSuccess: v.StealChangeSuccess,
					KillChangeSuccess:  v.KillChangeSuccess,
				}
				value = append(value, new)
			}

			myMonsters := model.MyMonsters{Monsters: value}
			objectIds := []*runtime.StorageWrite{
				{
					Collection:      "Monsters",
					Key:             "MyMonsters",
					UserID:          p.OwnerID,
					Value:           string(lo.Must1(json.Marshal(myMonsters))),
					PermissionWrite: 1,
					PermissionRead:  2,
				},
			}

			if _, err := nk.StorageWrite(ctx, objectIds); err != nil {
				logger.WithField("err", err).Error("Storage write error.")
			}
		}
	}

	return killed, nil
}
