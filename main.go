package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/rnikrozoft/hellomymonsters-backend/modules"
)

// noinspection GoUnusedExportedFunction
func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	// if err := initializer.RegisterRpc(rpcIdCanClaimDailyReward, modules.RpcCanClaimDailyReward); err != nil {
	// 	return err
	// }

	if err := initializer.RegisterRpc("steal_monster", modules.RpcStealMonster); err != nil {
		return err
	}

	if err := initializer.RegisterRpc("register", modules.RpcUserRegister); err != nil {
		logger.Error(err.Error())
		return err
	}

	// if err := initializer.RegisterRpc("kill_monster", modules.RpcKillMonster); err != nil {
	// 	logger.Error(err.Error())
	// 	return err
	// }

	logger.Info("Plugin loaded in '%d' msec.", time.Since(initStart).Milliseconds())
	return nil
}
