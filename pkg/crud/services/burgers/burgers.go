package burgers

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"webform/pkg/crud/models"
)

type BurgersSvc struct {
	pool *pgxpool.Pool
}

func NewBurgersSvc(pool *pgxpool.Pool) *BurgersSvc {
	if pool == nil {
		panic(errors.New("pool can't be nil")) // <- be accurate
	}
	return &BurgersSvc{pool: pool}
}

func (service *BurgersSvc) BurgersList() (list []models.Burger, err error) {
	list = make([]models.Burger, 0) // TODO: for REST API
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return nil, err // TODO: wrap to specific error
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(), searchWearFalse)
	if err != nil {
		return nil, err // TODO: wrap to specific error
	}
	defer rows.Close()

	for rows.Next() {
		item := models.Burger{}
		err := rows.Scan(&item.Id, &item.Name, &item.Price)
		if err != nil {
			return nil, err // TODO: wrap to specific error
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (service *BurgersSvc) Save(model models.Burger) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return err // TODO: wrap to specific error
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), insertToTable, model.Name, model.Price)
	if err != nil {
		return err
	}
	return nil
}

func (service *BurgersSvc) RemoveById(id int) (err error) {
	conn, err := service.pool.Acquire(context.Background())
	if err != nil {
		return err // TODO: wrap to specific error
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), updateInTable, id)
	if err != nil {
		return err
	}
	return nil
}
