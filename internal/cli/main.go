package cli

import (
	"context"
	"sync"

	"github.com/alecthomas/kingpin"
	"github.com/rarimo/voting-relayer/internal/config"
	"github.com/rarimo/voting-relayer/internal/data/pg"
	"github.com/rarimo/voting-relayer/internal/service"
	ingester "github.com/rarimo/voting-relayer/internal/state_transitor"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3"
)

func Run(args []string) bool {
	log := logan.New()

	defer func() {
		if rvr := recover(); rvr != nil {
			log.WithRecover(rvr).Error("app panicked")
		}
	}()

	cfg := config.New(kv.MustFromEnv())
	log = cfg.Log()

	app := kingpin.New("voting-relayer", "")

	runCmd := app.Command("run", "run command")
	serviceCmd := runCmd.Command("service", "run service") // you can insert custom help
	votingCmd := runCmd.Command("voting", "run voting")
	relayerCmd := runCmd.Command("relayer", "run relayer")

	migrateCmd := app.Command("migrate", "migrate command")
	migrateUpCmd := migrateCmd.Command("up", "migrate db up")
	migrateDownCmd := migrateCmd.Command("down", "migrate db down")

	// custom commands go here...

	cmd, err := app.Parse(args[1:])
	if err != nil {
		log.WithError(err).Error("failed to parse arguments")
		return false
	}

	switch cmd {
	case serviceCmd.FullCommand():
		ws := new(sync.WaitGroup)
		ws.Add(2)
		go service.Run(cfg, ws)
		go ingester.NewService(cfg, ingester.NewPassportRootIngester(cfg, pg.NewStateQ(cfg.DB())), ws).Run(context.Background())
		ws.Wait()

	case votingCmd.FullCommand():
		service.Run(cfg, nil)

	case relayerCmd.FullCommand():
		ingester.NewService(cfg, ingester.NewPassportRootIngester(cfg, pg.NewStateQ(cfg.DB())), nil).Run(context.Background())

	case migrateUpCmd.FullCommand():
		if err = MigrateUp(cfg); err != nil {
			log.WithError(err).Error("failed to migrate up")
			return false
		}

	case migrateDownCmd.FullCommand():
		if err = MigrateDown(cfg); err != nil {
			log.WithError(err).Error("failed to migrate down")
			return false
		}

	default:
		log.Errorf("unknown command %s", cmd)
		return false
	}

	return true
}
