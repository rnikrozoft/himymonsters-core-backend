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

	rpcRegistrations := map[string]func(context.Context, runtime.Logger, *sql.DB, runtime.NakamaModule, string) (string, error){
		"register":           modules.RpcUserRegister,
		"get_monsters":       modules.RpcGetMonsters,
		"kill_monster":       modules.RpcKillMonster,
		"free_gachapon":      modules.RpcFreeGachapon,
		"cheat_add_monsters": modules.RpcCheatAddMonsters,
		"steal_monster":      modules.RpcStealMonster,
	}

	for name, handler := range rpcRegistrations {
		if err := initializer.RegisterRpc(name, handler); err != nil {
			logger.Error(err.Error())
			return err
		}
	}

	logger.Info("Plugin loaded in '%d' msec.", time.Since(initStart).Milliseconds())
	return nil
}
