package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"

	"github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/rnikrozoft/hellomymonsters-backend/constant"
	"github.com/rnikrozoft/hellomymonsters-backend/model"
	"github.com/rnikrozoft/hellomymonsters-backend/utility"
	"github.com/samber/lo"
)

func RpcFreeGachapon(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("errNoUserIdFound")
	}

	records, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: "Gachapon",
			Key:        "Free",
			UserID:     userID,
		},
	})
	if err != nil {
		logger.WithField("err", err).Error("Storage read error.")
		return "", errors.New("not found free gachapon")
	}

	gachapon := model.Gachapon{}
	json.Unmarshal([]byte(records[0].Value), &gachapon)

	if !gachapon.IsCanCaim {
		return utility.Response(gachapon), nil
	}

	monsterList := []model.Monster{
		{
			Name:        "fallen_1",
			MonsterType: "fallen_1",
			Steal: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
			Kill: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
		},
		{
			Name:        "golem_3",
			MonsterType: "golem_3",
			Steal: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
			Kill: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
		},
		{
			Name:        "minotaur_3",
			MonsterType: "minotaur_3",
			Steal: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
			Kill: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
		},
		{
			Name:        "reaperman_1",
			MonsterType: "reaperman_1",
			Steal: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
			Kill: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
		},
		{
			Name:        "reaperman_2",
			MonsterType: "reaperman_2",
			Steal: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
			Kill: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
		},
		{
			Name:        "reaperman_3",
			MonsterType: "reaperman_3",
			Steal: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
			Kill: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
		},
	}

	got := model.MyMonsters{
		Monsters: make([]model.Monster, 0),
	}
	for i := 0; i < 2; i++ {
		random := rand.Intn(len(monsterList))
		monster := monsterList[random]
		monster.ID = uuid.New()
		got.Monsters = append(got.Monsters, monster)
	}

	objectIds := []*runtime.StorageWrite{
		{
			Collection: constant.Collection_Monsters,
			Key:        constant.Key_MyMonsters,
			UserID:     userID,
			Value:      string(lo.Must1(json.Marshal(got))),
		},
		{
			Collection: "Gachapon",
			Key:        "Free",
			UserID:     userID,
			Value: string(
				lo.Must1(
					json.Marshal(model.Gachapon{
						IsCanCaim: false,
					}),
				),
			),
		},
	}

	if _, err := nk.StorageWrite(ctx, objectIds); err != nil {
		logger.WithField("err", err).Error("Storage write error.")
		return "", err
	}

	return utility.Response(got), nil
}
