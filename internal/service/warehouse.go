package service

import (
	"database/sql"
	"errors"
	"la-warehouse/internal/model"

	"github.com/rs/zerolog"
)

type DB interface {
	CheckWarehouse(whID int64) error
	Reservation(whID int64, productID int64, amount int64) error
	ReservationFree(id int64, productID int64) error
	GetAmountProducts(id int64) (int64, error)
}

type Warehouse struct {
	log *zerolog.Logger
	pg  DB
}

func NewWarehouse(log *zerolog.Logger, pg DB) *Warehouse {
	return &Warehouse{log: log, pg: pg}
}

// Reservation резервирование товара на складе для доставки
func (w *Warehouse) Reservation(whID int64, products []model.ProductsAmount) (model.RespReservation, error) {
	err := w.pg.CheckWarehouse(whID)
	if err != nil {
		return model.RespReservation{}, err
	}

	resp := model.RespReservation{
		Products: make([]model.ProductsStatus, 0, len(products)),
	}
	for _, product := range products {
		status := "ok"
		if product.Amount <= 0 {
			resp.Products = append(resp.Products, model.ProductsStatus{
				ProductID: product.ProductID,
				Status:    "amount must be greater 0",
			})
			continue
		}

		err = w.pg.Reservation(whID, product.ProductID, product.Amount)
		if err != nil {
			w.log.Err(err).
				Int64("warehouse_id", whID).
				Int64("product_id", product.ProductID).
				Int64("amount", product.Amount).
				Msg("reservation failed")
			status = err.Error()
		}

		resp.Products = append(resp.Products, model.ProductsStatus{
			ProductID: product.ProductID,
			Status:    status,
		})
	}

	return resp, nil
}

// ReservationFree освобождение резерва товаров
func (w *Warehouse) ReservationFree(whID int64, products []int64) (model.RespReservationFree, error) {
	err := w.pg.CheckWarehouse(whID)
	if err != nil {
		return model.RespReservationFree{}, err
	}

	resp := model.RespReservationFree{
		Products: make([]model.ProductsStatus, 0, len(products)),
	}
	for _, productID := range products {
		status := "ok"
		err := w.pg.ReservationFree(whID, productID)
		if err != nil {
			w.log.Err(err).
				Int64("warehouse_id", whID).
				Int64("product_id", productID).
				Msg("reservation free failed")
			status = err.Error()
		}
		resp.Products = append(resp.Products, model.ProductsStatus{
			ProductID: productID,
			Status:    status,
		})
	}

	return resp, nil
}

// GetAmountProducts получение кол-ва оставшихся товаров на складе
func (w *Warehouse) GetAmountProducts(whID int64) (int64, error) {
	amount, err := w.pg.GetAmountProducts(whID)
	if errors.Is(err, sql.ErrNoRows) {
		amount = 0
	}
	return amount, nil
}
