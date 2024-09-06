package modules

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/rnikrozoft/hellomymonsters-backend/model"
	"github.com/samber/lo"
)

func RpcUserRegister(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("errNoUserIdFound")
	}

	var p *model.RegisterWithDeviceID
	if err := json.Unmarshal([]byte(payload), &p); err != nil {
		logger.Error(err.Error())
		return "", err
	}

	if p.Created {
		storageWrites := []*runtime.StorageWrite{
			{
				Collection: "Gachapon",
				Key:        "Free",
				UserID:     userID,
				Value: string(
					lo.Must1(
						json.Marshal(model.Gachapon{
							IsCanCaim: true,
						}),
					),
				),
			},
			{
				Collection: "Settings",
				Key:        "PlayTutorials",
				UserID:     userID,
				Value: string(
					lo.Must1(
						json.Marshal(map[string]bool{
							"done": false,
						}),
					),
				),
			},
		}
		walletUpdates := []*runtime.WalletUpdate{
			{
				UserID: userID,
				Changeset: map[string]int64{
					"coin":    1000,
					"dimonds": 100,
				},
			},
		}

		storageAcks, walletUpdateResults, err := nk.MultiUpdate(ctx, nil, storageWrites, nil, walletUpdates, false)
		if err != nil {
			logger.WithField("err", err).Error("Multi update error.")
		} else {
			logger.Info("Storage Acks: %d", len(storageAcks))
			logger.Info("Wallet Updates: %d", len(walletUpdateResults))
		}
	}
	return payload, nil
}
