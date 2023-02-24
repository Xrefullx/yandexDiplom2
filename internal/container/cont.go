package container

import (
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"github.com/Xrefullx/yandexDiplom2/internal/pkg"
)

func GetStorage() pkg.Storage {
	return Container.Get("storage").(pkg.Storage)
}

func GetConfig() models.Config {
	return Container.Get("server-config").(models.Config)
}
