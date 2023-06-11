package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"la-warehouse/internal/model"
)

// Reservation оборачивам резервацию в одну транзакцию.
// Блокируем запись в таблице products
func (p *Pg) Reservation(whID int64, productID int64, amount int64) (err error) {
	var tx *sql.Tx
	tx, err = p.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = errors.Join(err, tx.Commit())
		}
	}()

	result := tx.QueryRow(`SELECT id, name, size, amount 
		FROM public.products where id = $1 for update`, productID)
	if result.Err() != nil {
		return result.Err()
	}
	var product model.Product
	err = result.Scan(
		&product.ID,
		&product.Name,
		&product.Size,
		&product.Amount)
	if err != nil {
		return handleNotFound("products", err)
	}

	if amount > product.Amount {
		return errors.New("not enough amount of product")
	}

	_, err = tx.Exec(`UPDATE public.products set amount = $1 where id = $2`, product.Amount-amount, productID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO public.warehouse_product 
    		(product_id, warehouse_id, amount) values ($1, $2, $3) 
    		ON CONFLICT (product_id, warehouse_id) DO UPDATE
    		SET amount = warehouse_product.amount + $3`, productID, whID, amount)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pg) ReservationFree(id int64, productID int64) (err error) {
	var tx *sql.Tx
	tx, err = p.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = errors.Join(err, tx.Commit())
		}
	}()

	result := tx.QueryRow(`SELECT id, name, size, amount 
		FROM public.products where id = $1 for update`, productID)
	if result.Err() != nil {
		return result.Err()
	}
	var product model.Product
	err = result.Scan(
		&product.ID,
		&product.Name,
		&product.Size,
		&product.Amount)
	if err != nil {
		return handleNotFound("products", err)
	}

	row := tx.QueryRow(`DELETE FROM public.warehouse_product 
       WHERE product_id = $1 and warehouse_id = $2 
       RETURNING amount`, productID, id)
	if row.Err() != nil {
		return row.Err()
	}
	var amount int64
	err = row.Scan(&amount)
	if err != nil {
		return handleNotFound("warehouse_product", err)
	}

	_, err = tx.Exec(`UPDATE public.products set amount = $1 where id = $2`, product.Amount+amount, productID)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pg) GetAmountProducts(id int64) (int64, error) {
	row := p.QueryRow(`SELECT SUM(amount) from public.warehouse_product where warehouse_id = $1`, id)
	if row.Err() != nil {
		return 0, row.Err()
	}
	var amount int64
	err := row.Scan(&amount)
	if err != nil {
		return 0, handleNotFound("warehouse_product", err)
	}

	return amount, nil
}

// CheckWarehouse проверяем что склад существует и доступен
// TODO если запись в таблице склада будет меняться, то продумать блокировку склада
func (p *Pg) CheckWarehouse(whID int64) error {
	result := p.QueryRow(`SELECT available FROM public.warehouse WHERE id = $1`, whID)
	if result.Err() != nil {
		return result.Err()
	}
	var available bool
	err := result.Scan(&available)
	if err != nil {
		return handleNotFound("warehouse", err)
	}

	if !available {
		return errors.New("warehouse not available")
	}

	return nil
}

func handleNotFound(entityName string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%s not exist", entityName)
	}
	return err
}
