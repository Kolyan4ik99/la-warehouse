package model

// Warehouse сущность для описания склада
type Warehouse struct {
	// уникальный индетификатор
	ID int64
	// Name название
	Name string
	// IsAvailable признак доступности
	IsAvailable bool
}
