package main

import (
	"github.com/sullyh7/myportfolio/env"
	"github.com/sullyh7/myportfolio/internal/db"
	"github.com/sullyh7/myportfolio/internal/server"
	"github.com/sullyh7/myportfolio/internal/store"
	"go.uber.org/zap"
)

func main() {
	cfg := server.Config{
		Addr: env.GetString("ADDR", ":3000"),
		Db: server.DBConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/myportfolio?sslmode=disable"),
			MaxOpenConns: env.GetInt("MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("MAX_OPEN_CONNS", 30),
			MaxIdleTime:  env.GetString("MAX_IDE_TIME", "15m"),
		},
		Env:     env.GetString("ENV", "development"),
		Version: "1",
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(cfg.Db.Addr, cfg.Db.MaxOpenConns, cfg.Db.MaxIdleConns, cfg.Db.MaxIdleTime)
	if err != nil {
		logger.Fatalw("error connecting to db", "err", err)
	}
	defer db.Close()
	store := store.NewStorage(db)

	server := &server.Server{
		Config: cfg,
		Store:  store,
		Logger: logger,
	}
	mux := server.Mount()
	if err = server.Run(mux); err != nil {
		logger.Fatalw(err.Error())
	}
}
