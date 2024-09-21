package modules

import (
	"context"
	"database/sql"
	"errors"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/rnikrozoft/himymonsters-core-backend/constant"
)

func RpcGetShop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	_, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", errors.New("errNoUserIdFound")
	}

	objects, err := nk.StorageRead(ctx, []*runtime.StorageRead{{
		Collection: constant.Collection_Shop,
		Key:        constant.Key_Items,
		UserID:     constant.SystemID,
	}})
	if err != nil {
		logger.Error("StorageRead error: %v", err)
		return "", err
	}

	return objects[0].GetValue(), nil
}
