package modules

import (
	"context"
	"database/sql"
	"errors"

	"github.com/heroiclabs/nakama-common/runtime"
)

func RpcFreeGachapon(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("errNoUserIdFound")
	}

	objectIds := []*runtime.StorageRead{
		{
			Collection: "Gachapon",
			Key:        "Free",
			UserID:     userID,
		},
	}

	records, err := nk.StorageRead(ctx, objectIds)
	if err != nil {
		logger.WithField("err", err).Error("Storage read error.")
		return "", errors.New("not found free gachapon")
	}

	for _, record := range records {
		logger.Info("read: %d, write: %d, value: %s", record.PermissionRead, record.PermissionWrite, record.Value)
	}
	return payload, nil
}
