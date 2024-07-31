package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

type dailyReward struct {
	LastClaimUnix int64 `json:"last_claim_unix"`
}

func getLastDailyRewardObject(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, payload string) (dailyReward, *api.StorageObject, error) {
	var d dailyReward
	d.LastClaimUnix = 0

	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return d, nil, errors.New("errNoUserIdFound")
	}

	if len(payload) > 0 {
		return d, nil, errors.New("err")
	}

	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{{
		Collection: "reward",
		Key:        "daily",
		UserID:     userID,
	}})
	if err != nil {
		logger.Error("StorageRead error: %v", err)
		return d, nil, err
	}

	var o *api.StorageObject
	for _, object := range objects {
		switch object.GetKey() {
		case "daily":
			if err := json.Unmarshal([]byte(object.GetValue()), &d); err != nil {
				logger.Error("Unmarshal error: %v", err)
				return d, nil, err
			}
			return d, object, nil
		}
	}

	return d, o, nil
}

func canUserClaimDailyReward(d dailyReward) bool {
	t := time.Now()
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return time.Unix(d.LastClaimUnix, 0).Before(midnight)
}

func RpcCanClaimDailyReward(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var resp struct {
		CanClaimDailyReward bool `json:"canClaimDailyReward"`
	}

	dailyReward, _, err := getLastDailyRewardObject(ctx, logger, nk, payload)
	if err != nil {
		logger.Error("Error getting daily reward: %v", err)
		return "", err
	}

	resp.CanClaimDailyReward = canUserClaimDailyReward(dailyReward)

	out, err := json.Marshal(resp)
	if err != nil {
		logger.Error("Marshal error: %v", err)
		return "", err
	}

	logger.Debug("rpcCanClaimDailyReward resp: %v", string(out))
	return string(out), nil
}

func RpcClaimDailyReward(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("error not found")
	}

	var resp struct {
		CoinsReceived int64 `json:"coinsReceived"`
	}
	resp.CoinsReceived = int64(0)

	dailyReward, dailyRewardObject, err := getLastDailyRewardObject(ctx, logger, nk, payload)
	if err != nil {
		logger.Error("Error getting daily reward: %v", err)
		return "", err
	}

	if canUserClaimDailyReward(dailyReward) {
		resp.CoinsReceived = 500

		// Update player wallet.
		changeset := map[string]int64{
			"coins": resp.CoinsReceived,
		}
		if _, _, err := nk.WalletUpdate(ctx, userID, changeset, map[string]interface{}{}, false); err != nil {
			logger.Error("WalletUpdate error: %v", err)
			return "", err
		}

		err := nk.NotificationsSend(ctx, []*runtime.NotificationSend{{
			Code: 1001,
			Content: map[string]interface{}{
				"coins": changeset["coins"],
			},
			Persistent: true,
			Sender:     "", // Server sent.
			Subject:    "You've received your daily reward!",
			UserID:     userID,
		}})
		if err != nil {
			logger.Error("NotificationsSend error: %v", err)
			return "", err
		}

		dailyReward.LastClaimUnix = time.Now().Unix()

		object, err := json.Marshal(dailyReward)
		if err != nil {
			logger.Error("Marshal error: %v", err)
			return "", err
		}

		version := ""
		if dailyRewardObject != nil {
			// Use OCC to prevent concurrent writes.
			version = dailyRewardObject.GetVersion()
		}

		// Update daily reward storage object for user.
		_, err = nk.StorageWrite(ctx, []*runtime.StorageWrite{{
			Collection:      "reward",
			Key:             "daily",
			PermissionRead:  1,
			PermissionWrite: 0, // No client write.
			Value:           string(object),
			Version:         version,
			UserID:          userID,
		}})
		if err != nil {
			logger.Error("StorageWrite error: %v", err)
			return "", err
		}
	}

	out, err := json.Marshal(resp)
	if err != nil {
		logger.Error("Marshal error: %v", err)
		return "", err
	}

	logger.Debug("rpcClaimDailyReward resp: %v", string(out))
	return string(out), nil
}
