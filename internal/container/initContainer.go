package container

import (
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/Xrefullx/yandexDiplom2/internal/pkg/pg"
	"github.com/sarulabs/di"
	"go.uber.org/zap"
)

var Container di.Container

func Build(cfg models.Config, logger *zap.Logger) error {
	builder, err := di.NewBuilder()
	if err != nil {
		return err
	}
	var Storage pg.Storage
	if cfg.DataBaseURI != "" {
		_, _ = pg.New(cfg.DataBaseURI)
	}
	if err = Storage.Ping(); err != nil {
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
		Build: func(ctn di.Container) (interface{}, error) { return Storage, nil }}); err != nil {
		return err
	}
	Container = builder.Build()
	return nil
}
