package transport

import (
	"la-warehouse/internal/model"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json"
)

type WarehouseService interface {
	Reservation(whID int64, products []model.ProductsAmount) (model.RespReservation, error)
	ReservationFree(whID int64, products []int64) (model.RespReservationFree, error)
	GetAmountProducts(whID int64) (int64, error)
}

type Warehouse struct {
	service WarehouseService
}

func NewWarehouse(service WarehouseService) *Warehouse {
	return &Warehouse{service: service}
}

func (w *Warehouse) Reservation(r *http.Request, args *model.ReqReservation, reply *model.RespReservation) error {
	resp, err := w.service.Reservation(args.WhID, args.Products)
	if err != nil {
		return err
	}
	*reply = resp
	return nil
}

func (w *Warehouse) ReservationFree(r *http.Request, args *model.ReqReservationFree, reply *model.RespReservationFree) error {
	resp, err := w.service.ReservationFree(args.WhID, args.Products)
	if err != nil {
		return err
	}
	*reply = resp
	return nil
}

func (w *Warehouse) GetAmountProducts(r *http.Request, args *model.ReqGetAmountProducts, reply *model.RespGetAmountProducts) error {
	resp, err := w.service.GetAmountProducts(args.WhID)
	if err != nil {
		return err
	}
	reply.Amount = resp
	return nil
}

func (w *Warehouse) ListenAndServe(addr string) error {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	err := server.RegisterService(w, "")
	if err != nil {
		return err
	}
	r := mux.NewRouter()
	r.Handle("/rpc", server)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		return err
	}
	return nil
}
