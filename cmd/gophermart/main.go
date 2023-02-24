package main

import (
	"context"
	"flag"
	"github.com/Xrefullx/yandexDiplom2/internal/api/consta"
	"github.com/Xrefullx/yandexDiplom2/internal/api/container"
	"github.com/Xrefullx/yandexDiplom2/internal/api/handlers"
	"github.com/Xrefullx/yandexDiplom2/internal/api/server"
	"github.com/Xrefullx/yandexDiplom2/internal/api/service"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
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
		log.Fatalln("ошибка считывания конфига", zap.Error(err))
	}
	flag.Parse()
	if cfg.ReleaseMOD {
		zapLogger, err = zap.NewProduction()
	} else {
		zapLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalln(err)
	}
	zapLogger.Info("считана следующая конфигурация",
		zap.String("AddressServer", cfg.Address),
		zap.String("AccrualAddress", cfg.AccrualAddress),
		zap.Bool("ReleaseMOD", cfg.ReleaseMOD),
	)
	zapLogger.Debug("полная конфигурация", zap.Any("config", cfg))
	if err = container.ContainerBuild(cfg, zapLogger); err != nil {
		zapLogger.Fatal("ошибка запуска Di контейнера", zap.Error(err))
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
			err = service.CalculationLoyaltyPoints(ctx)
			if err != nil {
				zapLogger.Error("ошибка в работе модуля", zap.Error(err))
			}
		}
	}()
	r := handlers.Router(cfg)
	server.InitServer(r, cfg)
}
