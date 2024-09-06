package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/mroth/weightedrand/v2"
	"github.com/rnikrozoft/hellomymonsters-backend/constant"
	"github.com/rnikrozoft/hellomymonsters-backend/model"
	"github.com/rnikrozoft/hellomymonsters-backend/utility"
	"github.com/samber/lo"
)

func RpcCheatAddMonsters(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("errNoUserIdFound")
	}

	type cheatAddMonsters struct {
		Amount int `json:"amount"`
	}
	var req cheatAddMonsters
	json.Unmarshal([]byte(payload), &req)

	myMonster := model.MyMonsters{
		Monsters: make([]model.Monster, 0),
	}

	for i := 0; i < req.Amount; i++ {
		myMonster.Monsters = append(myMonster.Monsters, model.Monster{
			ID:                 uuid.New(),
			Name:               "monster",
			MonsterType:        "ork_1",
			StealChangeSuccess: rand.Intn(10),
			KillChangeSuccess:  rand.Intn(10),
		})
	}

	objectIds := []*runtime.StorageWrite{
		{
			Collection: constant.Collection_Monsters,
			Key:        constant.Key_MyMonsters,
			UserID:     userID,
			Value: string(lo.Must1(
				json.Marshal(myMonster),
			)),
		},
	}

	_, err := nk.StorageWrite(ctx, objectIds)
	if err != nil {
		logger.WithField("err", err).Error("Storage write error.")
		return "", err
	}
	return payload, nil
}
func RpcGetMonsters(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	type owner struct {
		OwnerID string `json:"owner_id"`
	}
	var req owner
	json.Unmarshal([]byte(payload), &req)

	objectIds := []*runtime.StorageRead{
		{
			Collection: constant.Collection_Monsters,
			Key:        constant.Key_MyMonsters,
			UserID:     req.OwnerID,
		},
	}

	res := &model.MyMonsters{
		Monsters: make([]model.Monster, 0),
	}
	records, err := nk.StorageRead(ctx, objectIds)
	if err != nil {
		logger.WithField("err", err).Error("Storage read error.")
		return utility.Response(res), err
	}

	if len(records) == 0 {
		return utility.Response(res), nil
	}

	if err := json.Unmarshal([]byte(records[0].Value), res); err != nil {
		logger.WithField("err", err).Error("JSON unmarshal error.")
		return utility.Response(res), err
	}

	return utility.Response(res), nil
}

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

func RpcStealMonster(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {

	var p *model.StealOrkill
	json.Unmarshal([]byte(payload), &p)

	objectIds := []*runtime.StorageRead{
		{
			Collection: "Monsters",
			Key:        "MyMonsters",
			UserID:     p.OwnerID,
		},
	}

	steal := "false"
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
		max := 10 - monster.StealChangeSuccess
		if max <= 0 {
			delete(monsters, p.MonsterID)
			steal = "true"
		} else {
			chooser, _ := weightedrand.NewChooser(
				weightedrand.NewChoice(true, monster.KillChangeSuccess),
				weightedrand.NewChoice(false, max),
			)
			result := chooser.Pick()

			if result {
				delete(monsters, p.MonsterID)
				steal = "true"
			}
		}

		value := []model.Monster{}
		if steal == "true" {
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
			objectFriend := []*runtime.StorageWrite{
				{
					Collection:      "Monsters",
					Key:             "MyMonsters",
					UserID:          p.OwnerID,
					Value:           string(lo.Must1(json.Marshal(myMonsters))),
					PermissionWrite: 1,
					PermissionRead:  2,
				},
			}

			if _, err := nk.StorageWrite(ctx, objectFriend); err != nil {
				logger.WithField("err", err).Error("Storage write error.")
			}

			userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
			if !ok {
				return "", errors.New("errNoUserIdFound")
			}

			objectIds := []*runtime.StorageRead{
				{
					Collection: constant.Collection_Monsters,
					Key:        constant.Key_MyMonsters,
					UserID:     userID,
				},
			}

			res := &model.MyMonsters{
				Monsters: make([]model.Monster, 0),
			}

			records, err := nk.StorageRead(ctx, objectIds)
			if err != nil {
				logger.WithField("err", err).Error("Storage read error.")
				return utility.Response(res), err
			}

			if err := json.Unmarshal([]byte(records[0].Value), res); err != nil {
				logger.WithField("err", err).Error("JSON unmarshal error.")
				return utility.Response(res), err
			}

			res.Monsters = append(res.Monsters, monster)

			objectUpdate := []*runtime.StorageWrite{
				{
					Collection:      constant.Collection_Monsters,
					Key:             constant.Key_MyMonsters,
					UserID:          userID,
					Value:           string(lo.Must1(json.Marshal(res))),
					PermissionWrite: 1,
					PermissionRead:  2,
				},
			}

			if _, err := nk.StorageWrite(ctx, objectUpdate); err != nil {
				logger.WithField("err", err).Error("Storage write error.")
			}
		}
	}

	return steal, nil
}
