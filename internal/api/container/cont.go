package container

import (
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/Xrefullx/yandexDiplom2/internal/storage"
	"github.com/Xrefullx/yandexDiplom2/internal/storage/pg"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
)

var DiContainer di.Container

func ContainerBuild(cfg models.Config, logger *zap.Logger) error {
	builder, err := di.NewBuilder()
	if err != nil {
		return err
	}
	var LoyalityStorage storage.LoyalityStorage
	if cfg.DataBaseURI != "" {
		LoyalityStorage, err = pg.New(cfg.DataBaseURI)
		if err != nil {
			return err
		}
	}
	if err = LoyalityStorage.Ping(); err != nil {
		return err
	}
	if err = builder.Add(di.Def{
		Name:  "server-config",
		Build: func(ctn di.Container) (interface{}, error) { return cfg, nil }}); err != nil {
		return err
	}
	if err = builder.Add(di.Def{
		Name:  "zap-logger",
		Build: func(ctn di.Container) (interface{}, error) { return logger, nil }}); err != nil {
		return err
	}
	if err = builder.Add(di.Def{
		Name:  "storage",
		Build: func(ctn di.Container) (interface{}, error) { return LoyalityStorage, nil }}); err != nil {
		return err
	}
	DiContainer = builder.Build()
	return nil
}
