package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"

	"github.com/google/uuid"
	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/mroth/weightedrand/v2"
	"github.com/rnikrozoft/himymonsters-core-backend/constant"
	"github.com/rnikrozoft/himymonsters-core-backend/model"
	"github.com/rnikrozoft/himymonsters-core-backend/utility"
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
			ID:          uuid.New(),
			Name:        "monster",
			MonsterType: "fallen_1",
			Steal: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
			Kill: model.StealOrKillSettings{
				RateSuccess: rand.Intn(10),
				Price:       int64(rand.Intn(10)),
			},
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
	res := &model.MyMonsters{
		Monsters: make([]model.Monster, 0),
	}

	type owner struct {
		OwnerID string `json:"owner_id"`
	}
	var req owner
	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		logger.WithField("err", err).Error(err.Error())
		return utility.Response(res), err
	}

	objectIds := []*runtime.StorageRead{
		{
			Collection: constant.Collection_Monsters,
			Key:        constant.Key_MyMonsters,
			UserID:     req.OwnerID,
		},
	}

	records, err := nk.StorageRead(ctx, objectIds)
	if err != nil {
		logger.WithField("err", err).Error(err.Error())
		return utility.Response(res), err
	}

	if len(records) == 0 {
		return utility.Response(res), nil
	}

	return records[0].GetValue(), nil
}

func RpcKillMonster(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	killed := "false"

	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return killed, errors.New("errNoUserIdFound")
	}

	account, err := nk.AccountGetId(ctx, userID)
	if err != nil {
		logger.WithField("err", err).Error("Get accounts error.")
		return killed, err
	}

	w := &model.Wallet{}
	if err := json.Unmarshal([]byte(account.Wallet), w); err != nil {
		logger.WithField("err", err).Error("cannot unmarshal")
		return killed, err
	}

	var p *model.StealOrkill
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		logger.WithField("err", err).Error("cannot unmarshal")
		return killed, err
	}

	records, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: constant.Collection_Monsters,
			Key:        constant.Key_MyMonsters,
			UserID:     p.OwnerID,
		},
	})
	if err != nil {
		logger.WithField("err", err).Error("Storage read error.")
		return killed, err
	}

	if len(records) == 0 {
		logger.WithField("err", err).Error("no found records")
		return killed, errors.New("errNoRecordFound")
	}

	var record model.Record
	if err := json.Unmarshal([]byte(records[0].Value), &record); err != nil {
		logger.WithField("err", err).Error("cannot unmarshal")
		return killed, err
	}

	if len(record.Monsters) == 0 {
		logger.Error("not found records")
		return killed, errors.New("errNoRecord.Monsters")
	}

	monsters := lo.Associate(record.Monsters, func(f model.Monster) (uuid.UUID, model.Monster) {
		return f.ID, f
	})

	monster := monsters[p.MonsterID]

	if w.Coin < monster.Kill.Price {
		return killed, nil //coin not enougth
	}

	max := 10 - monster.Kill.RateSuccess
	if max <= 0 {
		delete(monsters, p.MonsterID)
		killed = "true"
	} else {
		chooser, _ := weightedrand.NewChooser(
			weightedrand.NewChoice(true, monster.Kill.RateSuccess),
			weightedrand.NewChoice(false, max),
		)
		result := chooser.Pick()

		if result {
			delete(monsters, p.MonsterID)
			killed = "true"
		}
	}

	changeset := map[string]int64{
		"coin": -monster.Kill.Price,
	}
	metadata := map[string]interface{}{
		"kill_success": killed,
	}
	if _, _, err := nk.WalletUpdate(ctx, userID, changeset, metadata, true); err != nil {
		logger.WithField("err", err).Error("Wallet update error.")
		return killed, err
	}

	if killed == "true" {
		monsters := lo.MapToSlice(monsters, func(k uuid.UUID, v model.Monster) model.Monster {
			return v
		})

		friendMonsters := model.MyMonsters{
			Monsters: monsters,
		}

		new := []*runtime.StorageWrite{
			{
				Collection:      constant.Collection_Monsters,
				Key:             constant.Key_MyMonsters,
				UserID:          p.OwnerID,
				Value:           string(lo.Must1(json.Marshal(friendMonsters))),
				PermissionWrite: 1,
				PermissionRead:  2,
			},
		}

		if _, err := nk.StorageWrite(ctx, new); err != nil {
			logger.WithField("err", err).Error("Storage write error.")
			return "false", err
		}
	}
	return killed, nil
}

func RpcStealMonster(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	steal := "false"

	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return steal, errors.New("errNoUserIdFound")
	}

	account, err := nk.AccountGetId(ctx, userID)
	if err != nil {
		logger.WithField("err", err).Error("Get accounts error.")
		return steal, err
	}

	w := &model.Wallet{}
	if err := json.Unmarshal([]byte(account.Wallet), w); err != nil {
		logger.WithField("err", err).Error("cannot unmarshal")
		return steal, err
	}

	var p *model.StealOrkill
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		logger.WithField("err", err).Error(err.Error())
		return steal, err
	}

	records, err := nk.StorageRead(ctx, []*runtime.StorageRead{
		{
			Collection: constant.Collection_Monsters,
			Key:        constant.Key_MyMonsters,
			UserID:     p.OwnerID,
		},
	})
	if err != nil {
		logger.WithField("err", err).Error(err.Error())
		return steal, err
	}

	if len(records) == 0 {
		return steal, errors.New("not found any records")
	}

	var record model.Record
	if err := json.Unmarshal([]byte(records[0].Value), &record); err != nil {
		logger.WithField("err", err).Error(err.Error())
		return steal, err
	}

	if len(record.Monsters) == 0 {
		return steal, errors.New("not found any monsters")
	}

	monsters := lo.Associate(record.Monsters, func(f model.Monster) (uuid.UUID, model.Monster) {
		return f.ID, f
	})

	monster := monsters[p.MonsterID]

	if w.Coin < monster.Steal.Price {
		return steal, nil //coin not enougth
	}

	max := 10 - monster.Steal.RateSuccess
	if max <= 0 {
		delete(monsters, p.MonsterID)
		steal = "true"
	} else {
		chooser, _ := weightedrand.NewChooser(
			weightedrand.NewChoice(true, monster.Steal.RateSuccess),
			weightedrand.NewChoice(false, max),
		)
		result := chooser.Pick()

		if result {
			delete(monsters, p.MonsterID)
			steal = "true"
		}
	}

	changeset := map[string]int64{
		"coin": -monster.Steal.Price,
	}
	metadata := map[string]interface{}{
		"steal_success": steal,
	}
	if _, _, err := nk.WalletUpdate(ctx, userID, changeset, metadata, true); err != nil {
		logger.WithField("err", err).Error(err.Error())
		return steal, err
	}

	if steal == "true" {
		monsters := lo.MapToSlice(monsters, func(k uuid.UUID, v model.Monster) model.Monster {
			return v
		})

		friendMonsters := model.MyMonsters{
			Monsters: monsters,
		}

		userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
		if !ok {
			return "false", errors.New("errNoUserIdFound")
		}

		myMonsters := &model.MyMonsters{
			Monsters: make([]model.Monster, 0),
		}

		records, err := nk.StorageRead(ctx, []*runtime.StorageRead{
			{
				Collection: constant.Collection_Monsters,
				Key:        constant.Key_MyMonsters,
				UserID:     userID,
			},
		})
		if err != nil {
			logger.WithField("err", err).Error(err.Error())
			return utility.Response(myMonsters), err
		}

		if err := json.Unmarshal([]byte(records[0].Value), myMonsters); err != nil {
			logger.WithField("err", err).Error("JSON unmarshal error.")
			return utility.Response(myMonsters), err
		}

		myMonsters.Monsters = append(myMonsters.Monsters, monster)

		if _, err := nk.StorageWrite(ctx, []*runtime.StorageWrite{
			{
				Collection:      constant.Collection_Monsters,
				Key:             constant.Key_MyMonsters,
				UserID:          p.OwnerID,
				Value:           string(lo.Must1(json.Marshal(friendMonsters))),
				PermissionWrite: 1,
				PermissionRead:  2,
			},
			{
				Collection:      constant.Collection_Monsters,
				Key:             constant.Key_MyMonsters,
				UserID:          userID,
				Value:           string(lo.Must1(json.Marshal(myMonsters))),
				PermissionWrite: 1,
				PermissionRead:  2,
			},
		}); err != nil {
			logger.WithField("err", err).Error("Storage write error.")
			return steal, err
		}
	}

	return steal, nil
}
