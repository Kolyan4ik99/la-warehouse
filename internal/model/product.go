package model

// Product для описания сущности товар
type Product struct {
	// ID уникальный индентификатор
	ID int64 `db:"id"`
	// Name название
	Name string `db:"name"`
	// Size размер
	Size int64 `db:"size"`
	// Amount количество
	Amount int64 `db:"amount"`
}

type ProductsAmount struct {
	ProductID int64 `json:"id"`
	Amount    int64 `json:"amount"`
}

type ProductsStatus struct {
	ProductID int64  `json:"id"`
	Status    string `json:"status"`
}

type ReqReservation struct {
	WhID     int64            `json:"warehouseId"`
	Products []ProductsAmount `json:"products"`
}

type ReqReservationFree struct {
	WhID     int64   `json:"warehouseId"`
	Products []int64 `json:"products"`
}

type ReqGetAmountProducts struct {
	WhID int64 `json:"warehouseId"`
}

type RespReservation struct {
	Products []ProductsStatus `json:"products"`
}

type RespReservationFree struct {
	Products []ProductsStatus `json:"products"`
}

type RespGetAmountProducts struct {
	Amount int64 `json:"amount"`
}
