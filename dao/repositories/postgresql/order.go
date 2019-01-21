package postgresql

import (
	"database/sql"
	"github.com/pkg/errors"
	"github.com/netology/godesignpatterns/dao/models"
	"github.com/netology/godesignpatterns/dao/repositories"
)

func NewOrderRepository(db *sql.DB, orderItemRepository repositories.OrderItemRepository) repositories.OrderRepository {
	return &order{
		db:                  db,
		orderItemRepository: orderItemRepository,
	}
}

type order struct {
	db                  *sql.DB
	orderItemRepository repositories.OrderItemRepository
}

func (o *order) GetByID(orderID models.OrderID) (*models.Order, error) {
	stmt, err := o.db.Prepare("SELECT order_id, customer_id, amount, currency FROM orders WHERE order_id=$1")
	if err != nil {
		return nil, errors.Wrap(err, "prepare")
	}

	order := &models.Order{}
	err = stmt.QueryRow(orderID).Scan(&order.ID, &order.CustomerID, &order.Amount.Value, &order.Amount.Currency)
	if err != nil {
		return nil, errors.Wrap(err, "prepare")
	}

	orderItems, err := o.orderItemRepository.GetByOrderID(order.ID)
	if err != nil {
		return nil, errors.Wrap(err, "prepare")
	}
	order.Items = orderItems

	return order, nil
}

func (o *order) Save(order *models.Order) error {
	tx, err := o.db.Begin()
	if err != nil {
		return errors.Wrap(err, "begin transaction error")
	}

	stmt, err := tx.Prepare("INSERT INTO orders (customer_id, amount, currency) VALUES ($1, $2, $3) RETURNING order_id")
	if err != nil {
		return errors.Wrap(err, "prepare query error")
	}

	var lastInsertID int64
	if err := stmt.QueryRow(order.CustomerID, order.Amount.Value, order.Amount.Currency).Scan(&lastInsertID); err != nil {
		return errors.Wrap(err, "query row error")
	}

	/////////////// Alternative Usage: If pq (postgresql) driver support lastInsertID ////////////////////
	//
	// result, err := stmt.Exec(order.CustomerID, order.Amount.Value, order.Amount.Currency)
	// if err != nil {
	//	return errors.Wrap(err, "exec error")
	// }
	// lastInsertID, err = result.LastInsertId()
	//
	///////////////////////////////////////////////////////////////////////////////////////////////////////

	order.ID = models.OrderID(lastInsertID)
	for _, item := range order.Items {
		item.OrderID = order.ID
		if err := o.orderItemRepository.SaveWithTransaction(tx, item); err != nil {
			if e := tx.Rollback(); e != nil {
				return errors.Wrap(err, e.Error())
			}
			return errors.Wrap(err, "save order item error")
		}
	}

	if err = tx.Commit(); err != nil {
		return errors.Wrap(err, "commit error")
	}

	return nil
}
