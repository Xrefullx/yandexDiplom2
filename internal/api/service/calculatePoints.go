package service

import (
	"context"
	"encoding/json"
	"github.com/Xrefullx/yandexDiplom2/internal/api/consta"
	"github.com/Xrefullx/yandexDiplom2/internal/api/container"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"time"
)

func CalculationLoyaltyPoints(ctx context.Context) error {
	storage := container.GetStorage()
	cfg := container.GetConfig()
	log := container.GetLog()
	ctxStorage, cancelStorage := context.WithTimeout(ctx, consta.TimeOutStorage)
	orders, err := storage.GetOrdersProcess(ctxStorage)
	cancelStorage()
	if err != nil {
		return err
	}
	for _, order := range orders {
		ctx, cancel := context.WithTimeout(ctx, consta.TimeOutRequest)
		joinPath, errJoin := url.JoinPath(cfg.AccrualAddress, "/api/orders/", order.NumberOrder)
		if errJoin != nil {
			cancel()
			return errJoin
		}
		r, errGet := http.Get(joinPath)
		if errGet != nil {
			cancel()
			return errGet
		}
		switch r.StatusCode {
		case http.StatusTooManyRequests:
			log.Error("TooManyRequests")
			time.Sleep(consta.TimeSleepTooManyRequests)
		case http.StatusInternalServerError:
			body, errRead := io.ReadAll(r.Body)
			if errRead != nil {
				cancel()
				return errRead
			}
			log.Error("StatusInternalServerError", zap.String("body", string(body)))
			errClose := r.Body.Close()
			if errClose != nil {
				cancel()
				return errClose
			}
		case http.StatusOK:
			body, errRead := io.ReadAll(r.Body)
			if errRead != nil {
				cancel()
				return errRead
			}
			log.Debug("body", zap.String("body", string(body)))
			var loyaltyPoint models.Loyalty
			if errUnmarshal := json.Unmarshal(body, &loyaltyPoint); errUnmarshal != nil {
				cancel()
				return errUnmarshal
			}
			errClose := r.Body.Close()
			if errClose != nil {
				cancel()
				return errClose
			}
			if loyaltyPoint.Status == consta.OrderStatusREGISTERED {
				loyaltyPoint.Status = consta.OrderStatusPROCESSING
			}
			if order.Status != loyaltyPoint.Status {
				if err = storage.UpdateOrder(ctx, loyaltyPoint); err != nil {
					cancel()
					return err
				}
			}
		}
		cancel()
	}
	return nil
}
