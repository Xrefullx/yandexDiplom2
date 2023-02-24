package main

import (
	"context"
	"flag"
	"github.com/Xrefullx/yandexDiplom2/internal/api/router"
	"github.com/Xrefullx/yandexDiplom2/internal/api/server"
	"github.com/Xrefullx/yandexDiplom2/internal/container"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/Xrefullx/yandexDiplom2/internal/utils"
	"github.com/Xrefullx/yandexDiplom2/internal/utils/consta"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"log"
	"time"
)

var cfg models.Config

func init() {
	flag.StringVar(&cfg.Address, "a", cfg.Address, "the launch address of the HTTP server")
	flag.StringVar(&cfg.DataBaseURI, "d", cfg.DataBaseURI, "a string with the address of the database connection")
	flag.StringVar(&cfg.AccrualAddress, "r", cfg.AccrualAddress, "address of the accrual calculation system")
}
func main() {
	var zapLogger *zap.Logger
	var err error
	if err = env.Parse(&cfg); err != nil {
		log.Fatalln("config reading error", zap.Error(err))
	}
	flag.Parse()
	if err = container.Build(cfg, zapLogger); err != nil {
		zapLogger.Fatal("error launching the Di container", zap.Error(err))
	}
	defer func() {
		if err = container.GetStorage().Close(); err != nil {
			zapLogger.Fatal(consta.ErrorDataBase, zap.Error(err))
		}
	}()
	go func() {
		for {
			ctx := context.Background()
			time.Sleep(consta.TimeSleepCalculationLoyaltyPoints)
			err = utils.CalculationLoyaltyPoints(ctx)
			if err != nil {
				zapLogger.Error("ошибка ", zap.Error(err))
			}
		}
	}()
	r := router.Router()
	server.InitServer(r, cfg)
}
