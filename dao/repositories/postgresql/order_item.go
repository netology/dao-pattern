package postgresql

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/netology/godesignpatterns/dao/models"
	"github.com/netology/godesignpatterns/dao/repositories"
)

func NewOrderItemRepository(db *sql.DB) repositories.OrderItemRepository {
	return &orderItem{
		db: db,
	}
}

type orderItem struct {
	db *sql.DB
}

func (o orderItem) GetByOrderID(orderID models.OrderID) ([]*models.OrderItem, error) {
	stmt, err := o.db.Prepare("SELECT order_item_id, order_id, product_id, quantity, price, currency FROM order_items WHERE order_id=$1")
	if err != nil {
		return nil, errors.Wrap(err, "prepare")
	}

	rows, err := stmt.Query(orderID)
	if err != nil {
		return nil, errors.Wrap(err, "prepare")
	}

	orderItems := []*models.OrderItem{}
	for rows.Next() {
		orderItem := &models.OrderItem{}
		err = rows.Scan(&orderItem.ID, &orderItem.OrderID, &orderItem.ProductID, &orderItem.Quantity, &orderItem.Price.Value, &orderItem.Price.Currency)
		if err != nil {
			return nil, errors.Wrap(err, "prepare")
		}
		orderItems = append(orderItems, orderItem)
	}

	return orderItems, nil
}

func (o orderItem) Save(orderItem *models.OrderItem) error {
	stmt, err := o.db.Prepare("INSERT INTO order_items (order_id, product_id, quantity, price, currency) VALUES ($1, $2, $3, $4, $5) RETURNING order_item_id;")
	if err != nil {
		return err
	}

	var lastInsertID int64
	row := stmt.QueryRow(orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price.Value, orderItem.Price.Currency)
	if err := row.Scan(&lastInsertID); err != nil {
		return err
	}

	orderItem.ID = lastInsertID

	return nil
}

func (o orderItem) SaveWithTransaction(tx *sql.Tx, orderItem *models.OrderItem) error {
	stmt, err := tx.Prepare("INSERT INTO order_items (order_id, product_id, quantity, price, currency) VALUES ($1, $2, $3, $4, $5) RETURNING order_item_id;")
	if err != nil {
		return err
	}

	var lastInsertID int64
	row := stmt.QueryRow(orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price.Value, orderItem.Price.Currency)
	if err := row.Scan(&lastInsertID); err != nil {
		return err
	}

	orderItem.ID = lastInsertID

	return nil
}
